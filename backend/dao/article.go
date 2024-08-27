package dao

import (
	"blog/dao/db"
	"blog/utils/es"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// elasticsearch schema
const mapping = `{
			"settings": {
				"number_of_shards": 3,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"comma": {
							"type": "pattern",
							"pattern": ","
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"content": {
						"type": "text",
						"analyzer": "ik_max_word"
					},
					"tags": {
						"type": "text",
						"analyzer": "comma"
					},
					"title": {
						"type": "text",
						"analyzer": "ik_max_word"
					}
				}
			}
		}`

func GetArticle() *article {
	return articleDao
}

func init() {
	articleContentEsIndex := viper.GetString("article.contentsearchindex")
	db.GetMysql().AutoMigrate(&Article{})
	//init elasticsearch index and mapper
	es := db.GetElasticsearch()
	var err error
	var resp *esapi.Response
	resp, err = es.Indices.Exists([]string{articleContentEsIndex})
	if err != nil {
		logrus.Panic("check the index exist failed:", err.Error())
	}
	if resp.IsError() {
		resp, err = es.Indices.Create(articleContentEsIndex, es.Indices.Create.WithBody(strings.NewReader(mapping)))
		if err != nil {
			logrus.Panicf("create the index %s mapper %s failed: %s", articleContentEsIndex, mapping, err.Error())
		}
		if resp.IsError() {
			logrus.Panicf("create the index %s mapper %s failed: %s", articleContentEsIndex, mapping, resp.String())
		}
	}
	articleDao.searchEngine = es
	articleDao.esIndex = articleContentEsIndex
	articleDao.cachems = viper.GetInt64("cache.cleaninterval")
	articleDao.cachekeyPrefix = _a.TableName()
}

var articleContentBucketName string

type article struct {
	_              [0]func()
	searchEngine   *elasticsearch.Client
	esIndex        string
	cachems        int64
	cachekeyPrefix string
}

var articleDao = &article{}

func (a *article) CreateArticle(ctx context.Context, article *Article) (id uint, err error) {
	err = db.GetMysql().WithContext(ctx).Model(&Article{}).Create(article).Error
	if err != nil {
		return
	}
	id = article.ID
	abort := db.GetRedis().Set(ctx, fmt.Sprintf("%s_%d", a.cachekeyPrefix, article.ID), article, time.Millisecond*time.Duration(a.cachems)).Err()
	if abort != nil {
		logrus.Errorf("set article %v cache failed: %s", article, err.Error())
	}
	return
}

func (a *article) UpdateArticle(ctx context.Context, article *Article) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%d", a.cachekeyPrefix, article.ID)
	err = cache.Del(ctx, key).Err()
	if err != nil {
		logrus.Errorf("to update article,delete article from redis failed:%s", err.Error())
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&Article{}).Where("id = ?", article.ID).Updates(article).Error
	if err != nil {
		logrus.Errorf("update article %v from mysql failed:%s", article, err.Error())
		return
	}
	cache.Set(ctx, key, article, time.Millisecond*time.Duration(a.cachems))
	return
}
func (a *article) DeleteArticle(ctx context.Context, id uint) (err error) {
	err = db.GetRedis().Del(ctx, fmt.Sprintf("%s_%d", a.cachekeyPrefix, id)).Err()
	if err != nil {
		logrus.Errorf("delete the article %d cache failed: %s", id, err.Error())
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&Article{}).Where("id = ?", id).Delete(&Article{}).Error
	if err != nil {
		return
	}
	return
}
func (a *article) FindArticlePaticalById(ctx context.Context, id uint) (article Article, err error) {
	cache := db.GetRedis()
	err = cache.Get(ctx, fmt.Sprintf("%s_%d", a.cachekeyPrefix, id)).Scan(&article)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("find article %d from redis failed:%s", id, err.Error())
		}
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&Article{}).Select("id, title, creator, tags,created_at,images").Where("id = ?", id).First(&article).Error
	if err != nil {
		logrus.Errorf("find aritcle partical %d from mysql failed:%s", id, err.Error())
	}
	return
}
func (a *article) FindArticleById(ctx context.Context, id uint) (article Article, err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%d", a.cachekeyPrefix, id)
	err = cache.Get(ctx, key).Scan(&article)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("get article %d cache failed: %s", id, err.Error())
		}
		return
	}
	err = db.GetMysql().WithContext(ctx).Model(&Article{}).Where("id = ?", id).First(&article).Error
	if err != nil {
		logrus.Errorf("get article %d from mysql failed: %s", id, err.Error())
		return
	}
	ignoreErr := cache.Set(ctx, key, &article, 30*time.Minute).Err()
	if ignoreErr != nil {
		logrus.Errorf("set article %d cache failed: %s", id, ignoreErr.Error())
	}
	return
}

/*
*
通过文章内容，标签构建文章搜素引擎，用于文章搜索这里使用elasticsearch。

*
*/

func (a *article) BuildArticleSearch(ctx context.Context, article *Article) (err error) {
	req := esapi.CreateRequest{
		Index:      a.esIndex,
		DocumentID: strconv.Itoa(int(article.ID)),
	}
	//将文章内容、tags和标题放入req.body中
	var bd []byte
	bd, err = json.Marshal(struct {
		Content string `json:"content"`
		Tags    string `json:"tags"`
		Title   string `json:"title"`
	}{
		Content: strconv.Quote(article.Content),
		Tags:    article.Tags,
		Title:   article.Title,
	})
	req.Body = bytes.NewBuffer(bd)
	var resp *esapi.Response
	resp, err = req.Do(ctx, a.searchEngine)
	if err != nil {
		logrus.Error("build article search failed:", err.Error())
		return
	}
	if resp.IsError() {
		logrus.Error("build article search failed:", resp.String())
		err = &es.ESResponseError{}
		return
	}
	return
}

func (a *article) SearchArticleByPage(ctx context.Context, keyword string, page, size int) (articlesid []uint64, total uint, err error) {
	var req esapi.SearchRequest
	a.searchEngine.Search()
	req.Index = []string{a.esIndex}
	req.Body = strings.NewReader(`
		{
		   "query":{
			 "multi_match": {
			   "query": "` + keyword + `",
			   "fields": ["content^1","title^3","tags^10"]
			 }
		   },
		   "_source": false, 
		   "from":` + strconv.Itoa((page-1)*size) + `,
		   "size":` + strconv.Itoa(size) + `
		}
		`)
	var resp *esapi.Response
	resp, err = req.Do(ctx, a.searchEngine)
	if err != nil {
		logrus.Error("search article failed:", err.Error())
		return
	}
	if resp.IsError() {
		logrus.Error("search article failed:", resp.String())
		err = &es.ESResponseError{}
		return
	}
	var byt []byte
	byt, err = io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("read response body failed:", err.Error())
		return
	}
	resp.Body.Close()
	var re *es.SearchResult = new(es.SearchResult)
	err = json.Unmarshal(byt, re)
	if err != nil {
		logrus.Error("unmarshal response body failed:", err.Error())
		return
	}
	total = uint(re.Hits.TotalHits.Value)
	articlesid = make([]uint64, len(re.Hits.Hits))
	result := re.Hits.Hits
	for i, hit := range result {
		articlesid[i], err = strconv.ParseUint(hit.Id, 10, 64)
		if err != nil {
			logrus.Error("parse article id failed:", err.Error())
			return
		}
	}
	return
}

func (a *article) FindArticlePaticalByCreateTime(ctx context.Context, page, size int) (articles []*Article, total int64, err error) {
	storage := db.GetMysql()
	err = storage.WithContext(ctx).Model(&Article{}).Count(&total).Error
	if err != nil {
		return
	}
	err = storage.WithContext(ctx).Model(&Article{}).Select("id, title, creator, tags,created_at,images").Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&articles).Error
	if err != nil {
		return
	}
	return
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _a = &Article{}

/*文章表*/
type Article struct {
	gorm.Model
	Title   string `gorm:"type:varchar(255);not null"`
	Tags    string `gorm:"tags;varchar(300)"`
	Creator string `gorm:"varchar(64);not null"`
	Content string `gorm:"type:longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;not null"`
	Images  string `gorm:"type:longtext"`
}

func (a *Article) TableName() string {
	return "article"
}
func (a *Article) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Article) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
