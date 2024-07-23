package db

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	cfg := elasticsearch.Config{
		Addresses: viper.GetStringSlice("elasticsearch.address"),
		Username:  viper.GetString("elasticsearch.username"),
		Password:  viper.GetString("elasticsearch.password"),
		//TODO this method will stop the world
		// Transport: &fasthttp.Transport{},
	}
	var err error
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
}

var es *elasticsearch.Client

func GetElasticsearch() *elasticsearch.Client {
	return es
}
