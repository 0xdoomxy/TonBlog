package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	db.GetMysql().AutoMigrate(&Like{})
	likeDao = newLikeDao()
}

type like struct {
	_        [0]func()
	cacheKey string
	onceMaps map[uint]any
	rwmutex  sync.RWMutex
}

var likeDao *like

func newLikeDao() *like {
	return &like{
		cacheKey: viper.GetString("like.cachekeyPrefix"),
		onceMaps: make(map[uint]any),
		rwmutex:  sync.RWMutex{},
	}
}

func GetLike() *like {
	return likeDao
}

/*
*

	文章关注总览表

*
*/
type Like struct {
	ArticleID uint `gorm:"not null;uniqueIndex:search"`
	LikeNum   uint `gorm:"not null"`
}

func (like *Like) MarshalBinary() ([]byte, error) {
	return json.Marshal(like)
}

func (like *Like) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, like)
}

func (l *like) initLikeToCache(ctx context.Context, articleId uint) (err error) {
	storage := db.GetMysql()
	var like Like
	err = storage.Where("article_id = ?", articleId).First(&like).Error
	if err != nil {
		logrus.Errorf("init like to cache failed: %s", err.Error())
		return
	}
	cache := db.GetRedis()
	return cache.Set(ctx, fmt.Sprintf("%s_%d", l.cacheKey, articleId), like.LikeNum, 0).Err()
}

// before the execute the any like operation you should run this function,exclude the create like function
func (l *like) onceInitLikeToCache(ctx context.Context, articleid uint) (err error) {
	l.rwmutex.RLock()
	if _, ok := l.onceMaps[articleid]; !ok {
		l.rwmutex.RUnlock()
		l.rwmutex.Lock()
		l.onceMaps[articleid] = struct{}{}
		err = l.initLikeToCache(ctx, articleid)
		l.rwmutex.Unlock()
	} else {
		l.rwmutex.RUnlock()
	}
	return
}

func (l *like) IncrementLike(ctx context.Context, like *Like) (err error) {
	err = l.onceInitLikeToCache(ctx, like.ArticleID)
	if err != nil {
		logrus.Errorf("init like %v failed:%s", like, err.Error())
		return
	}
	cache := db.GetRedis()
	err = cache.IncrBy(ctx, fmt.Sprintf("%s_%d", l.cacheKey, like.ArticleID), int64(like.LikeNum)).Err()
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("increment like %d failed: %s", like.ArticleID, err.Error())
		}
		return
	}
	l.rwmutex.Lock()
	delete(l.onceMaps, like.ArticleID)
	l.rwmutex.Unlock()
	return l.IncrementLike(ctx, like)
}
func (l *like) DecrementLike(ctx context.Context, like *Like) (err error) {
	err = l.onceInitLikeToCache(ctx, like.ArticleID)
	if err != nil {
		logrus.Errorf("init like %v failed:%s", like, err.Error())
		return
	}
	cache := db.GetRedis()
	err = cache.DecrBy(ctx, fmt.Sprintf("%s_%d", l.cacheKey, like.ArticleID), int64(like.LikeNum)).Err()
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("decrement like %d failed: %s", like.ArticleID, err.Error())
		}
		return
	}
	l.rwmutex.Lock()
	delete(l.onceMaps, like.ArticleID)
	l.rwmutex.Unlock()
	return l.DecrementLike(ctx, like)

}

func (l *like) FindLikeById(ctx context.Context, articleid uint) (like Like, err error) {
	err = l.onceInitLikeToCache(ctx, articleid)
	if err != nil {
		logrus.Errorf("init like %v failed:%s", like, err.Error())
		return
	}
	like = Like{
		ArticleID: articleid,
	}
	cache := db.GetRedis()
	err = cache.Get(ctx, fmt.Sprintf("%s_%d", l.cacheKey, articleid)).Scan(&like.LikeNum)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("get like %d cache failed: %s", articleid, err.Error())
		}
		return
	}
	err = db.GetMysql().Where("article_id = ?", articleid).First(&like).Error
	if err != nil {
		logrus.Errorf("find like by id %d failed: %s ", articleid, err.Error())
	}
	l.compensateLike(articleid)
	return
}
func (l *like) DeleteLike(ctx context.Context, articleid uint) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("%s_%d", l.cacheKey, articleid)).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete the like %d cache failed: %s", articleid, err.Error())
		return
	}
	err = db.GetMysql().Where("article_id = ?", articleid).Delete(&Like{}).Error
	if err != nil {
		logrus.Errorf("delete the like %d failed: %s", articleid, err.Error())
	}
	l.compensateLike(articleid)

	return
}

func (l *like) CreateLike(ctx context.Context, like *Like) (err error) {
	err = db.GetMysql().Create(like).Error
	return
}

/*
*

	if the cache is removed, you should compensate the like num.

*
*/
func (l *like) compensateLike(articleid uint) {
	l.rwmutex.Lock()
	delete(l.onceMaps, articleid)
	l.rwmutex.Unlock()
}
