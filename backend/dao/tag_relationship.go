package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func GetTagRelationship() *tagRelationship {
	return tagRelationshipDao
}

type tagRelationship struct {
	cacheKey string
}

var tagRelationshipDao *tagRelationship = newTagRelationshipDao()

func newTagRelationshipDao() *tagRelationship {
	return &tagRelationship{
		cacheKey: _tr.TableName(),
	}
}

func (t *tagRelationship) CreateTagRelationship(ctx context.Context, tagRelationship *TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&TagRelationship{}).Create(tagRelationship).Error
	if err != nil {
		logrus.Errorf("create tag relationship %v failed:%v", tagRelationship, err)
		return
	}
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cacheKey, tagRelationship.Name)
	ignoreErr := cache.SAdd(ctx, key, tagRelationship.ArticleId).Err()
	if ignoreErr != nil && ignoreErr != redis.Nil {
		defer cache.Del(ctx, key)
		logrus.Errorf("add the tag relationship %v to redis failed:%s", tagRelationship, ignoreErr.Error())
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
	cache := db.GetRedis()
	var agg = make(map[string][]any)
	for _, tagRelationship := range tagRelationships {
		v, ok := agg[tagRelationship.Name]
		if !ok {
			v = make([]any, 0)
			agg[tagRelationship.Name] = v
		}
		v = append(v, tagRelationship.ArticleId)
	}
	for k, v := range agg {
		key := fmt.Sprintf("%s_%s", t.cacheKey, k)
		ignoreErr := cache.SAdd(ctx, key, v).Err()
		if ignoreErr != nil && ignoreErr != redis.Nil {
			defer cache.Del(ctx, key)
			logrus.Errorf("add the tag relationship %v to redis failed:%s", tagRelationships, ignoreErr.Error())
		}
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
	start := page - 1
	size := pagesize
	var articlesStr []string
	articlesStr, err = cache.SMembers(ctx, key).Result()
	if err != nil || len(articlesStr) > 0 {
		if err != nil {
			logrus.Errorf("find tag relationship %s from redis failed:%s", name, err.Error())
			return
		}
		size := max(size, len(articlesStr))
		view = make(TagRelationships, 0, pagesize)
		var articleid int
		var articleStr string
		if start >= len(articlesStr) {
			return
		}
		for i := start; i < size+start; i++ {
			if i >= len(articlesStr) {
				break
			}
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
	var ids = make([]any, 0, len(view))
	for _, v := range view {
		ids = append(ids, v.ArticleId)
	}
	ignoreErr := cache.SAdd(ctx, key, ids...).Err()
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

type TagRelationships []*TagRelationship

func (trs *TagRelationships) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, trs)
}
func (trs *TagRelationships) MarshalBinary() ([]byte, error) {
	return json.Marshal(trs)
}

func init() {
	db.GetMysql().AutoMigrate(&TagRelationship{})
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _tr = &TagRelationship{}

/*标签关系表*/
type TagRelationship struct {
	Name      string `gorm:"type:varchar(255);primary_key"`
	ArticleId uint   `gorm:"type:int;primary_key"`
}

func (tr *TagRelationship) TableName() string {
	return "tag_relationship"
}

func (tr *TagRelationship) MarshalBinary() ([]byte, error) {
	return json.Marshal(tr)
}
func (tr *TagRelationship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tr)
}
