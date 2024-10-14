package dao

import (
	"blog/dao/db"
	"blog/model"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func GetUser() *user {
	return userDao
}

type user struct {
	_        [0]func()
	cachekey string
	sf       singleflight.Group
}

var userDao *user = newUserDao()

func newUserDao() *user {
	return &user{
		cachekey: _u.TableName(),
		sf:       singleflight.Group{},
	}
}

func (u *user) CreateUser(user *model.User) (err error) {
	err = db.GetMysql().Model(&model.User{}).Create(user).Error
	if err != nil {
		return
	}
	return
}

func (u *user) FindUserByAddress(ctx context.Context, address string) (user model.User, err error) {
	var rawUser interface{}
	rawUser, err, _ = u.sf.Do(address, func() (inner_u interface{}, e error) {
		cache := db.GetRedis()
		key := fmt.Sprintf("%s_%s", u.cachekey, address)
		e = cache.Get(ctx, key).Scan(&inner_u)
		if !errors.Is(e, redis.Nil) {
			if e != nil {
				logrus.Errorf("find user %v failed from redis: %s", address, e.Error())
			}
			return
		}
		e = db.GetMysql().Model(&model.User{}).Where("address = ?", address).First(&inner_u).Error
		if e != nil {
			logrus.Errorf("find user %v failed from mysql:%s", address, e.Error())
			return
		}
		ignoreErr := cache.Set(ctx, key, &inner_u, 3*time.Minute).Err()
		if ignoreErr != nil {
			logrus.Errorf("set user %v to redis failed:%s", inner_u, ignoreErr.Error())
		}
		return
	})

	return rawUser.(model.User), err
}

func (u *user) DeleteUser(ctx context.Context, address string) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("%s_%s", u.cachekey, address)).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("delete user %v from redis failed:%s", address, err.Error())
		return
	}
	err = db.GetMysql().Model(&model.User{}).Where("address = ?", address).Delete(&model.User{}).Error
	if err != nil {
		logrus.Errorf("delete user %v failed from mysql:%s", address, err.Error())
	}
	return
}

func (u *user) UpdateUser(ctx context.Context, user *model.User) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", u.cachekey, user.Address)
	err = cache.Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Errorf("to update user, delete user from redis failed:%s", err.Error())
		return
	}
	err = db.GetMysql().Model(&model.User{}).Where("address = ?", user.Address).Updates(user).Error
	if err != nil {
		logrus.Errorf("update the user %v failed:%s", user, err.Error())
	}
	return
}

func init() {
	db.GetMysql().AutoMigrate(&model.User{})
}

// should replace the origin cacheKey which should assign the value by user. then we pass the tag table name to assign the cache prefix
var _u = &model.User{}
