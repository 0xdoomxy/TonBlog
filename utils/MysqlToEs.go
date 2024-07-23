package utils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("read config failed: %v", err)
	}
}

//es field

//	"content": {
//		"type": "text",
//		"analyzer": "ik_max_word"
//	},
//
//	"tags": {
//		"type": "text",
//		"analyzer": "comma"
//	},
//
//	"title": {
//		"type": "text",
//		"analyzer": "ik_max_word"
//	}
type Entity struct {
	ID      uint   `gorm:"id" json:"-"`
	Content string `gorm:"content" json:"content"`
	Title   string `gorm:"title" json:"title"`
	Tags    string `gorm:"tags" json:"tags"`
}

func main() {
	dbConn := newMysqlConn()
	esConn := newElasticsearchConn()
	var entities []*Entity
	var err error
	err = dbConn.Table(viper.GetString("mysql.table")).Scan(&entities).Error
	if err != nil {
		logrus.Panicf("query entities from mysql failed: %s", err.Error())
	}
	logrus.Infof("query entities from mysql total:%d", len(entities))
	var body []byte
	var request *esapi.IndexRequest
	var resp *esapi.Response
	for i := 0; i < len(entities); i++ {
		var e = entities[i]
		body, err = json.Marshal(e)
		if err != nil {
			logrus.Errorf("marshal entity failed: %s", err.Error())
		}
		request = &esapi.IndexRequest{
			Index:      viper.GetString("elasticsearch.index"),
			DocumentID: strconv.Itoa(int(e.ID)),
			Body:       bytes.NewReader(body),
			Refresh:    "true",
		}
		resp, err = request.Do(context.TODO(), esConn)
		if err != nil {
			logrus.Errorf("index entity failed: %s", err.Error())
		}
		defer resp.Body.Close()
		if resp.IsError() {
			logrus.Printf("[%s] Error indexing document ID=%d", resp.Status(), e.ID)
		} else {
			logrus.Printf("[%s] Successfully indexed document ID=%d", resp.Status(), e.ID)
		}
	}
}

func newMysqlConn() (db *gorm.DB) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", viper.GetString("mysql.username"), viper.GetString("mysql.password"), viper.GetString("mysql.host"), viper.GetInt("mysql.port"), viper.GetString("mysql.database"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		logrus.Panicf("connect to mysql %s failed: %s", dsn, err.Error())
	}
	var t *sql.DB
	t, err = db.DB()
	if err != nil {
		logrus.Panicf("failed to ping database: %s", err.Error())
	}
	err = t.Ping()
	if err != nil {
		logrus.Panicf("failed to ping database: %s", err.Error())
	}
	logrus.Info("connect to mysql success")
	return
}

func newElasticsearchConn() (es *elasticsearch.Client) {
	var err error
	cfg := elasticsearch.Config{
		Addresses: viper.GetStringSlice("elasticsearch.address"),
		Username:  viper.GetString("elasticsearch.username"),
		Password:  viper.GetString("elasticsearch.password"),
		//TODO this method will stop the world
		// Transport: &fasthttp.Transport{},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logrus.Panicf("connect to elasticsearch %v failed: %s", cfg, err.Error())
	}
	var resp *esapi.Response
	resp, err = es.Ping()
	if err != nil {
		logrus.Panicf("ping elasticsearch %v failed: %s", cfg, err.Error())
	}
	defer resp.Body.Close()
	logrus.Info("connect to elasticsearch success")
	return
}
