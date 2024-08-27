package cron

import (
	"blog/dao"
	"blog/dao/db"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type likeConsumerCron struct {
	internal *cron.Cron
}

func NewLikeConsumerCron() *likeConsumerCron {
	return &likeConsumerCron{
		internal: cron.New(),
	}
}

func (lcc *likeConsumerCron) Run() {
	var onceExec = make(chan struct{}, 1)
	var m = make(map[uint64]uint64)
	onceExec <- struct{}{}
	lcc.internal.AddJob("*/2 * * * *", cron.FuncJob(func() {
		if len(onceExec) < 1 {
			return
		}
		<-onceExec
		defer func() {
			onceExec <- struct{}{}
		}()
		var cache = db.GetRedis()
		var err error
		var keys []string
		cachekeys := dao.GetLikePreifxKey()
		keys, err = cache.Keys(context.TODO(), fmt.Sprintf("%s_*", cachekeys)).Result()
		if err != nil {
			logrus.Errorf("get keys from redis failed: %s", err.Error())
			return
		}
		var likenum uint64
		var articleidStr string
		var articleid uint64
		var found bool
		wg := sync.WaitGroup{}
		wg.Add(len(keys))
		for _, k := range keys {
			go func(key string) {
				defer wg.Done()
				var err error
				likenum, err = cache.Get(context.TODO(), key).Uint64()
				if err != nil {
					return
				}

				articleidStr, found = strings.CutPrefix(key, fmt.Sprintf("%s_", dao.GetLikePreifxKey()))
				if !found {
					return
				}
				articleid, err = strconv.ParseUint(articleidStr, 10, 64)
				if err != nil {
					logrus.Errorf("parse articleid error (articleid:%s,likenum:%d) failed: %s", articleidStr, likenum, err.Error())
					return
				}
				if old, ok := m[articleid]; !ok || (ok && old < likenum) {
					err = db.GetMysql().Model(&dao.Like{}).Where("article_id = ?", articleid).Update("like_num", likenum).Error
					if err != nil {
						logrus.Errorf("update like num (articleid:%d,likenum:%d) failed: %s", articleid, likenum, err.Error())
						return
					}
					m[articleid] = likenum
				} else if ok && old >= likenum {
					cache.Del(context.TODO(), key)
				}
			}(k)
		}
		wg.Wait()
	}))
	lcc.internal.Start()
}
