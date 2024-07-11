package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	db.GetMysql().AutoMigrate(&User{})
	userDao = newUserDao()
}

type user struct {
	_        [0]func()
	cachekey string
}

var userDao *user

func newUserDao() *user {
	return &user{
		cachekey: viper.GetString("user.cachekeyPrefix"),
	}
}
func GetUser() *user {
	return userDao
}

/*
*

	用户表

*
*/
type User struct {
	Address   string `gorm:"type:varchar(64);primary_key"`
	Alias     string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
}

func (user *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(user)
}

func (user *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
}

func (u *user) CreateUser(user *User) (err error) {
	err = db.GetMysql().Model(&User{}).Create(user).Error
	if err != nil {
		return
	}
	return
}

func (u *user) FindUserByAddress(ctx context.Context, address string) (user User, err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", u.cachekey, address)
	err = cache.Get(ctx, key).Scan(&user)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("find user %v failed from redis: %s", address, err.Error())
		}
		return
	}
	err = db.GetMysql().Model(&User{}).Where("address = ?", address).First(&user).Error
	if err != nil {
		logrus.Errorf("find user %v failed from mysql:%s", address, err.Error())
		return
	}
	ignoreErr := cache.Set(ctx, key, &user, 3*time.Minute).Err()
	if ignoreErr != nil {
		logrus.Errorf("set user %v to redis failed:%s", user, ignoreErr.Error())
	}
	return
}

func (u *user) DeleteUser(ctx context.Context, address string) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("%s_%s", u.cachekey, address)).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete user %v from redis failed:%s", address, err.Error())
		return
	}
	err = db.GetMysql().Model(&User{}).Where("address = ?", address).Delete(&User{}).Error
	if err != nil {
		logrus.Errorf("delete user %v failed from mysql:%s", address, err.Error())
	}
	return
}

func (u *user) UpdateUser(ctx context.Context, user *User) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%s", u.cachekey, user.Address)
	err = cache.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("to update user, delete user from redis failed:%s", err.Error())
		return
	}
	err = db.GetMysql().Model(&User{}).Where("address = ?", user.Address).Updates(user).Error
	if err != nil {
		logrus.Errorf("update the user %v failed:%s", user, err.Error())
	}
	return
}
