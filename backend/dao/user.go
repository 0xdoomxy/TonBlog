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
	"gorm.io/gorm"
)

func init() {
	db.GetMysql().AutoMigrate(&User{})
}

type user struct {
	cachekey string
}

var userDao *user = newUserDao()

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
	gorm.Model
	Address string `gorm:"type:varchar(256);not null"`
	Alias   string `gorm:"type:varchar(255);not null"`
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

func (u *user) FindUserById(ctx context.Context, userid uint) (user User, err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%d", u.cachekey, userid)
	err = cache.Get(ctx, key).Scan(&user)
	if err != redis.Nil {
		if err != nil {
			logrus.Errorf("find user %v failed from redis: %s", userid, err.Error())
		}
		return
	}
	err = db.GetMysql().Model(&User{}).Where("id = ?", userid).First(&user).Error
	if err != nil {
		logrus.Errorf("find user %v failed from mysql:%s", userid, err.Error())
	}
	cache.Set(ctx, key, &user, 3*time.Minute)
	return
}

func (u *user) DeleteUser(ctx context.Context, userid uint) (err error) {
	cache := db.GetRedis()
	err = cache.Del(ctx, fmt.Sprintf("%s_%d", u.cachekey, userid)).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("delete user %v from redis failed:%s", userid, err.Error())
		return
	}
	err = db.GetMysql().Model(&User{}).Where("id = ?", userid).Delete(&User{}).Error
	if err != nil {
		logrus.Errorf("delete user %v failed from mysql:%s", userid, err.Error())
	}
	return
}

func (u *user) UpdateUser(ctx context.Context, user *User) (err error) {
	cache := db.GetRedis()
	key := fmt.Sprintf("%s_%d", u.cachekey, user.ID)
	err = cache.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		logrus.Errorf("to update user, delete user from redis failed:%s", err.Error())
		return
	}
	err = db.GetMysql().Model(&User{}).Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		logrus.Errorf("update the user %v failed:%s", user, err.Error())
	}
	return
}
