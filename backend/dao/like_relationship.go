package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	//固定每个文章点赞的cache 个数,防止太多导致内存溢出
	likeScript = `
	local key = KEYS[1]
	local limit = tonumber(ARGV[1])
	local length = redis.call('ZCARD', key)
	local removed = {}
	if length > limit then
		removed = redis.call('ZRANGE', key, 0, length-limit-1)
    	redis.call('ZREMRANGEBYRANK', key, 0, length-limit-1)
	end
	return removed
	`
	//在程序结束时候执行的脚本用来将redis里面的数据存储到mysql
	finalizeScript = `
	local members = redis.call('ZRANGE', KEYS[1], 0, -1)
redis.call('DEL', KEYS[1])
return members`
)

func init() {
	db.GetMysql().AutoMigrate(&LikeRelationship{})
	likeRelationshipDao = newLikeRelationshipDao()
	go func() {
		/**
		监听程序退出信号，将redis里面的数据存储到mysql
		**/
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case <-sigs:
				cache := db.GetRedis()
				var res any
				var err error

				var userPublicKeys []any
				var ok bool
				for articleid := range likeRelationshipDao.times {
					res, err = cache.Eval(context.Background(), finalizeScript, []string{fmt.Sprintf("%s_%d", likeRelationshipDao.cacheKeyPrefix, articleid)}).Result()
					if err != nil {
						logrus.Errorf("dump the like relationship failed: %v", err)
						continue
					}
					userPublicKeys, ok = res.([]any)
					if !ok {
						logrus.Errorf("dump the like relationship failed: %v,type:%v", res, reflect.TypeOf(res))
						continue
					}
					if len(userPublicKeys) <= 0 {
						continue
					}
					storage := db.GetMysql()
					var models = make([]*LikeRelationship, len(userPublicKeys))
					for i, publickey := range userPublicKeys {
						models[i] = &LikeRelationship{
							ArticleID: articleid,
							PublicKey: publickey.(string),
						}
					}
					if err != nil {
						return
					}
					err = storage.Model(&likeRelationship{}).CreateInBatches(&models, (len(models)/100)+1).Error
					if err != nil {
						logrus.Errorf("dump the like relationship failed: %v", err)
					}
				}

			default:
				time.Sleep(time.Second * 5)
			}
		}
	}()
}

type likeRelationship struct {
	_              [0]func()
	cacheKeyPrefix string
	times          map[uint]*int32
	dumpFunc       sync.Pool
	maxcount       int32
	mutex          *sync.Mutex
}

func newLikeRelationshipDao() (res *likeRelationship) {
	maxcount := viper.GetInt32("like.relationship.maxcount")
	res = &likeRelationship{
		cacheKeyPrefix: viper.GetString("like.relationship.cachekeyPrefix"),
		times:          make(map[uint]*int32),
		maxcount:       maxcount,
		mutex:          &sync.Mutex{},
	}
	/**
	dumpFunc 用来将redis里面的数据存储到mysql，redis里面的数据超过maxcount的部分将会被存储到mysql
	**/
	dumpFunc := func(maxcount int32, articleid uint) {
		cache := db.GetRedis()
		res, err := cache.Eval(context.Background(), likeScript, []string{fmt.Sprintf("%s_%d", res.cacheKeyPrefix, articleid)}, maxcount).Result()
		if err != nil {
			if err != redis.Nil {
				logrus.Errorf("dump the like relationship failed: %v", err)
			}
			return
		}
		userPublickeys, ok := res.([]any)
		if !ok {
			logrus.Errorf("dump the like relationship failed: %v", res)
			return
		}
		storage := db.GetMysql()
		if len(userPublickeys) <= 0 {
			return
		}
		var models = make([]*LikeRelationship, len(userPublickeys))
		for i, publickey := range userPublickeys {
			models[i] = &LikeRelationship{
				ArticleID: articleid,
				PublicKey: publickey.(string),
			}
		}
		if err != nil {
			return
		}
		err = storage.Model(&LikeRelationship{}).CreateInBatches(&models, (len(models)/100)+1).Error
		if err != nil {
			logrus.Errorf("dump the like relationship %v failed: %v", models, err)
		}
	}
	res.dumpFunc = sync.Pool{
		New: func() any {
			return dumpFunc
		},
	}
	return
}

var likeRelationshipDao *likeRelationship

func GetLikeRelationship() *likeRelationship {
	return likeRelationshipDao
}

/*
文章关注表
*/
type LikeRelationship struct {
	ArticleID uint   `gorm:"primarykey"`
	PublicKey string `gorm:"varchar(64);primarykey"`
}

func (lrs *LikeRelationship) MarshalBinary() ([]byte, error) {
	return json.Marshal(lrs)
}

func (lrs *LikeRelationship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, lrs)
}

func (l *likeRelationship) CreateLikeRelationship(ctx context.Context, likeRelationship *LikeRelationship) (err error) {
	cache := db.GetRedis()
	err = cache.ZAdd(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), redis.Z{Score: float64(time.Now().Unix()), Member: likeRelationship.PublicKey}).Err()
	if err != nil {
		logrus.Errorf("create the like relationship %v failed: %v", likeRelationship, err)
		return
	}
	go func(articleid uint) {
		var old, ok = l.times[likeRelationship.ArticleID]
		if !ok {
			old = new(int32)
			l.times[likeRelationship.ArticleID] = old
		}
		/**
		避免锁竞争
		**/
		if *old < l.maxcount {
			*old++
			return
		}
		l.mutex.Lock()
		if *old < l.maxcount {
			*old++
			l.mutex.Unlock()
			return
		}
		go l.dumpFunc.Get().(func(int32, uint))(l.maxcount, articleid)
		*old = 0
		l.mutex.Unlock()
	}(likeRelationship.ArticleID)
	return
}

func (l *likeRelationship) DeleteLikeRelationship(ctx context.Context, likeRelationship *LikeRelationship) (err error) {
	cache := db.GetRedis()
	var res int64
	res, err = cache.ZRem(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), likeRelationship.PublicKey).Result()
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete the like relationship %v failed: %s", likeRelationship, err)
		return
	}
	if res <= 0 {
		// the like relationship is storaged by the mysql
		err = db.GetMysql().WithContext(ctx).Model(&LikeRelationship{}).Where("article_id = ? and public_key = ?", likeRelationship.ArticleID, likeRelationship.PublicKey).Delete(&LikeRelationship{}).Error
		if err != nil {
			logrus.Errorf("delete the like relationship %v failed: %s", likeRelationship, err.Error())
			return
		}
	}
	return
}

func (l *likeRelationship) FindLikeRelationshipByArticleID(ctx context.Context, likeRelationship *LikeRelationship) (likeRelationships []*LikeRelationship, err error) {
	cache := db.GetRedis()
	var userPublickeyStr []string
	userPublickeyStr, err = cache.ZRange(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), 0, -1).Result()
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
		}
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&LikeRelationship{}).Where("article_id = ?", likeRelationship.ArticleID).Find(&likeRelationships).Error
	if err != nil {
		logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
		return
	}
	var publickey string
	for i := 0; i < len(userPublickeyStr); i++ {
		publickey = userPublickeyStr[i]
		if err != nil {
			logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
			return
		}
		likeRelationships = append(likeRelationships, &LikeRelationship{
			ArticleID: likeRelationship.ArticleID,
			PublicKey: publickey,
		})
	}
	return
}

func (l *likeRelationship) FindLikeRelationshipByArticleIDAndUserid(ctx context.Context, likeRelationship *LikeRelationship) (exist bool, err error) {
	cache := db.GetRedis()
	_, err = cache.ZScore(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), likeRelationship.PublicKey).Result()
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
			return
		}
		exist = true
		return
	}
	var count int64
	err = db.GetMysql().WithContext(ctx).Model(&LikeRelationship{}).Where("article_id = ? and public_key = ? ", likeRelationship.ArticleID, likeRelationship.PublicKey).Count(&count).Error
	if err != nil {
		logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
		return
	}
	if count > 0 {
		exist = true
	}
	return
}
