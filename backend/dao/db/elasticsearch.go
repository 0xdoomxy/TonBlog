package db

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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
		logrus.Panic("connect to elasticsearch failed:", err.Error())
	}
	var resp *esapi.Response
	resp, err = es.Ping()
	if err != nil {
		logrus.Panic("ping elasticsearch failed:", err.Error())
	}
	defer resp.Body.Close()
	logrus.Info("connect to elasticsearch success")
}

var es *elasticsearch.Client

func GetElasticsearch() *elasticsearch.Client {
	return es
}
