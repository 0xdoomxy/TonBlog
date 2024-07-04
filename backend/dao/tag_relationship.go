package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	tagRelationshipDao = newTagRelationshipDao()
}

type TagRelationship struct {
	Name      string `gorm:"type:varchar(255);primary_key"`
	ArticleId uint   `gorm:"type:int;primary_key"`
}

func (tr *TagRelationship) MarshalBinary() ([]byte, error) {
	return json.Marshal(tr)
}
func (tr *TagRelationship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tr)
}

type TagRelationships []*TagRelationship

func (trs *TagRelationships) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, trs)
}
func (trs *TagRelationships) MarshalBinary() ([]byte, error) {
	return json.Marshal(trs)
}

type tagRelationship struct {
	cacheKey string
}

var tagRelationshipDao *tagRelationship

func newTagRelationshipDao() *tagRelationship {
	return &tagRelationship{
		cacheKey: viper.GetString("tag.relationship.cacheKeyPrefix"),
	}
}

func GetTagRelationship() *tagRelationship {
	return tagRelationshipDao
}

func (t *tagRelationship) CreateTagRelationship(ctx context.Context, tagRelationship *TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&TagRelationship{}).Create(tagRelationship).Error
	if err != nil {
		logrus.Errorf("create tag relationship %v failed:%v", tagRelationship, err)
		return
	}

	return
}

/*
*
批量创建 tag -article 关系
*/
func (t *tagRelationship) BatchCreateTagRelationship(ctx context.Context, tagRelationships []*TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&TagRelationship{}).Create(&tagRelationships).Error
	if err != nil {
		logrus.Errorf("batch create tag relationship %v failed:%v", tagRelationships, err)
	}
	return
}

/*
*
批量删除 tag -article 关系
*/
func (t *tagRelationship) BatchDeleteTagRelationship(ctx context.Context, tagRelationships []*TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&TagRelationship{}).Delete(&tagRelationships).Error
	if err != nil {
		logrus.Errorf("batch delete tag relationship %v failed:%v", tagRelationships, err)
	}
	var keym = make(map[string]any)
	for _, tagRelationship := range tagRelationships {
		key := fmt.Sprintf("%s_%s", t.cacheKey, tagRelationship.Name)
		keym[key] = struct{}{}
	}
	cache := db.GetRedis()
	for key, _ := range keym {
		ignoreErr := cache.Del(ctx, key).Err()
		if ignoreErr != nil {
			logrus.Errorf("Cache inconsistency:delete the key %s failed:%s", key, ignoreErr.Error())
		}
	}
	return
}
func (t *tagRelationship) FindTagRelationshipByName(ctx context.Context, name string, page int, pagesize int) (view TagRelationships, err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cacheKey, name)
	start := (pagesize - 1)
	size := pagesize
	var articlesStr []string
	articlesStr, err = cache.SMembers(ctx, key).Result()
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("find tag relationship %s from redis failed:%s", name, err.Error())
			return
		}
		view = make(TagRelationships, 0, pagesize)
		var articleid int
		var articleStr string
		for i := start; i < size; i++ {
			articleStr = articlesStr[i]
			articleid, err = strconv.Atoi(articleStr)
			if err != nil {
				logrus.Errorf("convert the articleid %s failed:%v", articleStr, err)
				return
			}
			view = append(view, &TagRelationship{
				Name:      name,
				ArticleId: uint(articleid),
			})
		}
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&TagRelationship{}).Where("name = ?", name).Limit(pagesize).Offset(start).Scan(&view).Error
	if err != nil {
		logrus.Errorf("find tag relationship %s failed:%v", name, err)
		return
	}
	var ids = make([]int, 0, len(view))
	for _, v := range view {
		ids = append(ids, int(v.ArticleId))
	}
	ignoreErr := cache.SAdd(ctx, key, ids).Err()
	if ignoreErr != nil {
		logrus.Errorf("add the tag relationship %v to redis failed:%s", view, ignoreErr.Error())
	}
	return
}
func (t *tagRelationship) DeleteTagRelationship(ctx context.Context, tagRelationship *TagRelationship) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cacheKey, tagRelationship.Name)
	err = cache.SRem(ctx, key, tagRelationship.ArticleId).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete the tag relationship %v failed:%s", tagRelationship, err.Error())
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&TagRelationship{}).Where("name = ? and article_id = ?", tagRelationship.Name, tagRelationship.ArticleId).Delete(&TagRelationship{}).Error
	return
}
