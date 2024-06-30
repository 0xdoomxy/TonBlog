package db

import (
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ossclient *oss.Client

var cache sync.Map = sync.Map{}

func init() {
	var err error
	ossclient, err = oss.New(viper.GetString("oss.endpoint"), viper.GetString("oss.accesskeyid"), viper.GetString("oss.accesskeysecret"))
	if err != nil {
		logrus.Fatal("阿里云 oss client 初始化失败", err.Error())
	}
	logrus.Info("阿里云 oss client 初始化成功")
}

func GetBucket(bucketName string) *oss.Bucket {
	if v, ok := cache.Load(bucketName); ok {
		return v.(*oss.Bucket)
	}
	bucket, err := ossclient.Bucket(bucketName)
	if err != nil {
		logrus.Error("获取 bucket 失败", err.Error())
		return nil
	}
	cache.Store(bucketName, bucket)
	return bucket
}

func GetOss() *oss.Client {
	return ossclient
}
