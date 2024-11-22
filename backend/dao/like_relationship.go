package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

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

func GetLikeRelationship() *likeRelationship {
	return likeRelationshipDao
}

func init() {
	db.GetMysql().AutoMigrate(&model.LikeRelationship{})
	likeRelationshipDao = newLikeRelationshipDao()
	go func() {
		/**
		监听程序退出信号，将redis里面的数据存储到mysql
		**/
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigs)
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
					var models = make([]*model.LikeRelationship, len(userPublicKeys))
					for i, publickey := range userPublicKeys {
						models[i] = &model.LikeRelationship{
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
				return
			default:
				time.Sleep(time.Second * 5)
			}
		}
	}()
}

var likeRelationshipDao *likeRelationship

type likeRelationship struct {
	_              [0]func()
	cacheKeyPrefix string
	times          map[uint]*int32
	dumpFunc       sync.Pool
	maxcount       int32
	mutex          *sync.Mutex
	sf             singleflight.Group
}

func newLikeRelationshipDao() (res *likeRelationship) {
	maxcount := viper.GetInt32("like.relationship.maxcount")
	res = &likeRelationship{
		cacheKeyPrefix: _lr.TableName(),
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
		var models = make([]*model.LikeRelationship, len(userPublickeys))
		for i, publickey := range userPublickeys {
			models[i] = &model.LikeRelationship{
				ArticleID: articleid,
				PublicKey: publickey.(string),
			}
		}
		if err != nil {
			return
		}
		err = storage.Model(&model.LikeRelationship{}).CreateInBatches(&models, (len(models)/100)+1).Error
		if err != nil {
			logrus.Errorf("dump the like relationship %v failed: %v", models, err)
		}
	}
	res.dumpFunc = sync.Pool{
		New: func() any {
			return dumpFunc
		},
	}
	res.sf = singleflight.Group{}
	return
}

func (l *likeRelationship) CreateLikeRelationship(ctx context.Context, likeRelationship *model.LikeRelationship) (err error) {
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

func (l *likeRelationship) DeleteLikeRelationship(ctx context.Context, likeRelationship *model.LikeRelationship) (err error) {
	cache := db.GetRedis()
	var isDelete bool = false

	defer func() {
		if err == nil && !isDelete {
			err = fmt.Errorf("delete like relationship cant exist")
		}
	}()
	var res int64
	res, err = cache.ZRem(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), likeRelationship.PublicKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("delete the like relationship %v failed: %s", likeRelationship, err)
		return
	}
	if err == nil {
		isDelete = true
	}
	if res <= 0 {
		// the like relationship is storaged by the mysql
		var deleteRes *gorm.DB
		deleteRes = db.GetMysql().WithContext(ctx).Model(&model.LikeRelationship{}).Where("article_id = ? and public_key = ?", likeRelationship.ArticleID, likeRelationship.PublicKey).Delete(&model.LikeRelationship{})
		if deleteRes.Error != nil {
			logrus.Errorf("delete the like relationship %v failed: %s", likeRelationship, deleteRes.Error.Error())
			return
		} else if deleteRes.RowsAffected > 0 {
			isDelete = true
		}
	}
	return
}

func (l *likeRelationship) FindLikeRelationshipByArticleID(ctx context.Context, likeRelationship *model.LikeRelationship) (likeRelationships []*model.LikeRelationship, err error) {

	cache := db.GetRedis()
	var userPublickeyStr []string
	userPublickeyStr, err = cache.ZRange(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), 0, -1).Result()
	if !errors.Is(err, redis.Nil) {
		if err != nil {
			logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
		}
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&model.LikeRelationship{}).Where("article_id = ?", likeRelationship.ArticleID).Find(&likeRelationships).Error
	if err != nil {
		logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, err)
		return
	}
	var publickey string
	for i := 0; i < len(userPublickeyStr); i++ {
		publickey = userPublickeyStr[i]
		likeRelationships = append(likeRelationships, &model.LikeRelationship{
			ArticleID: likeRelationship.ArticleID,
			PublicKey: publickey,
		})
	}
	return
}

func (l *likeRelationship) FindLikeRelationshipByArticleIDAndUserid(ctx context.Context, likeRelationship *model.LikeRelationship) (exist bool, err error) {
	var rawExist interface{}
	rawExist, err, _ = l.sf.Do(fmt.Sprintf("like_relationship_article_%d_user_%s", likeRelationship.ArticleID, likeRelationship.PublicKey), func() (inner_e interface{}, e error) {
		inner_e = false
		cache := db.GetRedis()
		_, e = cache.ZScore(ctx, fmt.Sprintf("%s_%d", l.cacheKeyPrefix, likeRelationship.ArticleID), likeRelationship.PublicKey).Result()
		if !errors.Is(e, redis.Nil) {
			if e != nil {
				logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, e)
				return
			}
			inner_e = true
			return
		}
		var count int64
		e = db.GetMysql().WithContext(ctx).Model(&model.LikeRelationship{}).Where("article_id = ? and public_key = ? ", likeRelationship.ArticleID, likeRelationship.PublicKey).Count(&count).Error
		if e != nil {
			logrus.Errorf("find the like relationship by articleid %v failed: %v", likeRelationship, e)
			return
		}
		if count > 0 {
			inner_e = true
		}
		return
	})

	return rawExist.(bool), err
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _lr = &model.LikeRelationship{}
