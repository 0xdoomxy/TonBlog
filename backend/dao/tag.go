package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// cant repeat from any tag cache key name
const ALL_TAGS_CACHE_KEY = "tags"

func GetTag() *tag {
	return tagDao
}

type tag struct {
	_        [0]func()
	cachekey string
	sf       singleflight.Group
}

var tagDao *tag = newTagDao()

func newTagDao() *tag {
	return &tag{
		cachekey: _t.TableName(),
		sf:       singleflight.Group{},
	}
}

func (t *tag) CreateTag(ctx context.Context, tag *model.Tag) (err error) {
	err = db.GetMysql().Model(&model.Tag{}).WithContext(ctx).Create(tag).Error
	if err != nil {
		return
	}
	ignoreErr := db.GetRedis().Del(ctx, ALL_TAGS_CACHE_KEY).Err()
	if ignoreErr != nil {
		logrus.Errorf("delete all tags from redis failed:%s", ignoreErr.Error())
	}
	return
}

func (t *tag) FindTag(ctx context.Context, name string) (tag model.Tag, err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cachekey, name)
	err = cache.Get(ctx, key).Scan(&tag)
	if !errors.Is(err, redis.Nil) {
		if err != nil {
			logrus.Errorf("get tag %s from redis failed:%v", name, err.Error())
		}
		return
	}
	err = db.GetMysql().Model(&model.Tag{}).WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		logrus.Errorf("get tag %s from mysql failed:%s", name, err.Error())
		return
	}
	cache.Set(ctx, key, tag, 0)
	return
}

func (t *tag) DeleteTag(ctx context.Context, name string) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cachekey, name)
	err = cache.Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("delete tag %s from redis failed:%s", name, err.Error())
		return
	}
	err = db.GetMysql().Model(&model.Tag{}).WithContext(ctx).Where("name = ?", name).Delete(&model.Tag{}).Error
	return
}

func (t *tag) FindAndIncrementTagNumByName(ctx context.Context, name string, num uint) (err error) {
	//乐观
	var needCreate bool = false
	cache := db.GetRedis()
	defer func() {
		if needCreate {
			err = db.GetMysql().Model(&model.Tag{}).Create(&model.Tag{Name: name, ArticleNum: num}).Error
			if err != nil {
				logrus.Errorf("create the tag (name:%s) failed:%s", name, err.Error())
				return
			}
			ignoreErr := cache.Del(ctx, ALL_TAGS_CACHE_KEY).Err()
			if ignoreErr != nil {
				logrus.Errorf("delete all tags from redis failed:%s", ignoreErr.Error())
			}
		}
	}()
	key := fmt.Sprintf("%s_%s", t.cachekey, name)
	err = cache.Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("delete tag %s from redis failed:%s", name, err.Error())
		return
	}
	result := db.GetMysql().WithContext(ctx).Model(&model.Tag{}).Where("name = ?", name).Update("article_num", gorm.Expr("article_num + ?", num))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		needCreate = true
	}
	return
}
func (t *tag) FindAllTags(ctx context.Context) (tags Tags, err error) {
	var rawTags interface{}
	rawTags, err, _ = t.sf.Do("all_tags", func() (interface{}, error) {
		var inner_t Tags
		var e error
		cache := db.GetRedis()
		e = cache.Get(ctx, ALL_TAGS_CACHE_KEY).Scan(&inner_t)
		if !errors.Is(e, redis.Nil) {
			if e != nil {
				logrus.Errorf("get all tags from redis failed:%s", e.Error())
			}
			return inner_t, e
		}
		e = db.GetMysql().Model(&model.Tag{}).WithContext(ctx).Find(&inner_t).Error
		if e != nil {
			//TODO if not records in database,the logic will error
			logrus.Errorf("get all tags from  mysql:%s", e.Error())
			return inner_t, e
		}
		ignoreErr := cache.Set(ctx, ALL_TAGS_CACHE_KEY, &inner_t, 2*time.Minute).Err()
		if ignoreErr != nil {
			logrus.Errorf("all tags set the redis error :%s", ignoreErr.Error())
		}
		return inner_t, e
	})
	return rawTags.(Tags), err
}

type Tags []*model.Tag

func (tags *Tags) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tags)
}
func (tags *Tags) MarshalBinary() ([]byte, error) {
	return json.Marshal(tags)
}

func init() {
	db.GetMysql().AutoMigrate(&model.Tag{})
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _t = &model.Tag{}
