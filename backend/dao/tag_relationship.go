package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func GetTagRelationship() *tagRelationship {
	return tagRelationshipDao
}

type tagRelationship struct {
	cacheKey string
	sf       singleflight.Group
}

var tagRelationshipDao *tagRelationship = newTagRelationshipDao()

func newTagRelationshipDao() *tagRelationship {
	return &tagRelationship{
		cacheKey: _tr.TableName(),
		sf:       singleflight.Group{},
	}
}

func (t *tagRelationship) CreateTagRelationship(ctx context.Context, tagRelationship *model.TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&model.TagRelationship{}).Create(tagRelationship).Error
	if err != nil {
		logrus.Errorf("create tag relationship %v failed:%v", tagRelationship, err)
		return
	}
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cacheKey, tagRelationship.Name)
	ignoreErr := cache.SAdd(ctx, key, tagRelationship.ArticleId).Err()
	if ignoreErr != nil && !errors.Is(ignoreErr, redis.Nil) {
		defer cache.Del(ctx, key)
		logrus.Errorf("add the tag relationship %v to redis failed:%s", tagRelationship, ignoreErr.Error())
	}
	return
}

/*
*
批量创建 tag -article 关系
*/
func (t *tagRelationship) BatchCreateTagRelationship(ctx context.Context, tagRelationships []*model.TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&model.TagRelationship{}).Create(&tagRelationships).Error
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
		if ignoreErr != nil && !errors.Is(ignoreErr, redis.Nil) {
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
func (t *tagRelationship) BatchDeleteTagRelationship(ctx context.Context, tagRelationships []*model.TagRelationship) (err error) {
	err = db.GetMysql().WithContext(ctx).Model(&model.TagRelationship{}).Delete(&tagRelationships).Error
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
	var rawTagRelationships interface{}
	rawTagRelationships, err, _ = t.sf.Do(fmt.Sprintf("tag_relationship_byname_%s_%d_%d", name, page, pagesize), func() (interface{}, error) {
		var inner_t TagRelationships = make(TagRelationships, 0, pagesize)
		var e error
		cache := db.GetRedis()
		key := fmt.Sprintf("%s_%s", t.cacheKey, name)
		start := page - 1
		size := pagesize
		var articlesStr []string
		articlesStr, e = cache.SMembers(ctx, key).Result()
		if e != nil || len(articlesStr) > 0 {
			if e != nil {
				logrus.Errorf("find tag relationship %s from redis failed:%s", name, e.Error())
				return inner_t, e
			}
			size := max(size, len(articlesStr))
			var articleid int
			var articleStr string
			if start >= len(articlesStr) {
				return inner_t, e
			}
			for i := start; i < size+start; i++ {
				if i >= len(articlesStr) {
					break
				}
				articleStr = articlesStr[i]
				articleid, e = strconv.Atoi(articleStr)
				if e != nil {
					logrus.Errorf("convert the articleid %s failed:%v", articleStr, e)
					return inner_t, e
				}
				inner_t = append(inner_t, &model.TagRelationship{
					Name:      name,
					ArticleId: uint(articleid),
				})
			}
			return inner_t, e
		}
		err = db.GetMysql().WithContext(ctx).Model(&model.TagRelationship{}).Where("name = ?", name).Limit(pagesize).Offset(start).Scan(&inner_t).Error
		if err != nil {
			logrus.Errorf("find tag relationship %s failed:%v", name, err)
			return inner_t, e
		}
		var ids = make([]any, 0, len(inner_t))
		for _, v := range inner_t {
			ids = append(ids, v.ArticleId)
		}
		ignoreErr := cache.SAdd(ctx, key, ids...).Err()
		if ignoreErr != nil {
			logrus.Errorf("add the tag relationship %v to redis failed:%s", inner_t, ignoreErr.Error())
		}
		return inner_t, e
	})

	return rawTagRelationships.(TagRelationships), err
}
func (t *tagRelationship) DeleteTagRelationship(ctx context.Context, tagRelationship *model.TagRelationship) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", t.cacheKey, tagRelationship.Name)
	err = cache.SRem(ctx, key, tagRelationship.ArticleId).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("delete the tag relationship %v failed:%s", tagRelationship, err.Error())
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&model.TagRelationship{}).Where("name = ? and article_id = ?", tagRelationship.Name, tagRelationship.ArticleId).Delete(&model.TagRelationship{}).Error
	return
}

type TagRelationships []*model.TagRelationship

func (trs *TagRelationships) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, trs)
}
func (trs *TagRelationships) MarshalBinary() ([]byte, error) {
	return json.Marshal(trs)
}

func init() {
	err := db.GetMysql().AutoMigrate(&model.TagRelationship{})
	if err != nil {
		logrus.Panicf("auto migrate tag_relationship table error:%s", err.Error())
	}
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _tr = &model.TagRelationship{}
