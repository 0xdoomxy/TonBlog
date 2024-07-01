package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// cant repeat from any tag cache key name
const ALL_TAGS_CACHE_KEY = "tags"

func init() {
	db.GetMysql().AutoMigrate(&Tag{})
}

type tag struct {
	_        [0]func()
	cachekey string
}

var tagDao *tag = newTagDao()

func newTagDao() *tag {
	return &tag{
		cachekey: viper.GetString("tag.cachekeyPrefix"),
	}
}
func GetTag() *tag {
	return tagDao
}

/*
*

	标签表

*
*/
type Tag struct {
	Name       string `gorm:"type:varchar(255);primaryKey"`
	ArticleNum uint   `gorm:"not null"`
}

func (tag *Tag) MarshalBinary() ([]byte, error) {
	return json.Marshal(tag)
}

func (tag *Tag) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tag)
}

func (t *tag) CreateTag(ctx context.Context, tag *Tag) (err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Create(tag).Error
	if err != nil {
		return
	}
	return
}

func (t *tag) FindTag(ctx context.Context, name string) (tag Tag, err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cachekey, name)
	err = cache.Get(ctx, key).Scan(&tag)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("get tag %s from redis failed:%v", name, err.Error())
		}
		return
	}
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("name = ?", name).First(&tag).Error
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
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete tag %s from redis failed:%s", name, err.Error())
		return
	}
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("name = ?", name).Delete(&Tag{}).Error
	return
}

func (t *tag) FindAndIncrementTagNumByName(ctx context.Context, name string, num uint) (err error) {
	//乐观
	var needCreate bool = false
	defer func() {
		if needCreate {
			err = db.GetMysql().Model(&Tag{}).Create(&Tag{Name: name, ArticleNum: num}).Error
			if err != nil {
				logrus.Errorf("create the tag (name:%s) failed:%s", name, err.Error())
			}
		}
	}()
	key := fmt.Sprintf("%s_%s", t.cachekey, name)
	cache := db.GetRedis()
	err = cache.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete tag %s from redis failed:%s", name, err.Error())
		return
	}
	result := db.GetMysql().WithContext(ctx).Model(&Tag{}).Where("name = ?", name).Update("article_num", gorm.Expr("article_num + ?", num))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		needCreate = true
	}
	return
}
func (t *tag) FindAllTags(ctx context.Context) (tags []*Tag, err error) {
	cache := db.GetRedis()
	err = cache.Get(ctx, ALL_TAGS_CACHE_KEY).Scan(&tags)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("get all tags from redis failed:%s", err.Error())
		}
		return
	}
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Find(&tags).Error
	if err != nil {
		logrus.Errorf("get all tags from  mysql:%s", err.Error())
		return
	}
	cache.Set(ctx, ALL_TAGS_CACHE_KEY, &tags, 0)
	return
}
