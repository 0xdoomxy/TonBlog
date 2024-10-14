package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"strconv"
)

func GetComment() *comment {
	return commentDao
}

func init() {
	db.GetMysql().AutoMigrate(&model.Comment{})
	commentDao = newCommentDao()
}

type comment struct {
	_        [0]func()
	cachekey string
	sf       singleflight.Group
}

var commentDao *comment

func newCommentDao() *comment {
	return &comment{
		cachekey: _c.TableName(),
		sf:       singleflight.Group{},
	}
}

func (c *comment) CreateComment(ctx context.Context, comment *model.Comment) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&model.Comment{}).Create(comment).Error
	if err != nil {
		logrus.Errorf("create the comment failed: %v", err)
		return
	}
	cache := db.GetRedis()
	ignoreErr := cache.HSet(ctx, fmt.Sprintf("%s_%d", c.cachekey, comment.ArticleID), strconv.Itoa(int(comment.ID)), comment).Err()
	if ignoreErr != nil {
		logrus.Errorf("create the comment cache failed: %v", err)
	}
	return
}

func (c *comment) FindCommentCreateBy(ctx context.Context, id uint, creator string) (ok bool, err error) {
	err = db.GetMysql().WithContext(ctx).Model(&model.Comment{}).Where("article_id = ? and creator = ?", id, creator).First(&model.Comment{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("find the comment by id %d and creator %s failed: %v", id, creator, err)
		}
		return
	}
	ok = true
	return
}

func (c *comment) DeleteComment(ctx context.Context, articleid uint, id uint) (err error) {
	cache := db.GetRedis()
	var del int64
	del, err = cache.HDel(ctx, fmt.Sprintf("%s_%d", c.cachekey, articleid), strconv.Itoa(int(id))).Result()
	if err != nil || del <= 0 {
		logrus.Errorf("delete the comment cache by articleid %d failed: %v", articleid, err)
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).Delete(&model.Comment{}).Error
	if err != nil {
		logrus.Errorf("delete the comment by id %d  failed: %v", id, err)
	}
	return
}

func (c *comment) DeleteCommentByArticle(ctx context.Context, articleid uint) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("%s_%d", c.cachekey, articleid)).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("delete the comment cache by articleid %d failed: %v", articleid, err)
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&model.Comment{}).Where("article_id = ?", articleid).Delete(&model.Comment{}).Error
	if err != nil {
		logrus.Errorf("delete the comment by articleid %d failed: %v", articleid, err)
	}
	return
}

func (c *comment) FindCommentByArticleid(ctx context.Context, articleid uint) (view []*model.Comment, err error) {
	var rawComments interface{}
	rawComments, err, _ = c.sf.Do(fmt.Sprintf("comment_article_%d", articleid), func() (interface{}, error) {
		inner_c := make([]*model.Comment, 0)
		var e error
		cache := db.GetRedis()
		if cache.Exists(ctx, fmt.Sprintf("%s_%d", c.cachekey, articleid)).Val() > 0 {
			e = cache.HVals(ctx, fmt.Sprintf("%s_%d", c.cachekey, articleid)).ScanSlice(&inner_c)
			if e != nil {
				logrus.Errorf("find the comment by articleid %d failed: %v", articleid, e)
			}
			return inner_c, e
		}
		e = db.GetMysql().WithContext(ctx).Model(&model.Comment{}).Where("article_id = ?", articleid).Find(&inner_c).Error
		if e != nil {
			logrus.Errorf("find the comment by articleid %d failed: %v", articleid, e)
		}
		var caches = make(map[string]interface{})
		for _, v := range inner_c {
			caches[strconv.Itoa(int(v.ID))] = v
		}
		ignoreErr := cache.HMSet(ctx, fmt.Sprintf("%s_%d", c.cachekey, articleid), caches).Err()
		if ignoreErr != nil {
			logrus.Errorf("set the comment cache by articleid %d failed: %v", articleid, ignoreErr)
		}
		return inner_c, e
	})
	return rawComments.([]*model.Comment), err
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _c = &model.Comment{}
