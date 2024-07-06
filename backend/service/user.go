package service

import (
	"blog/dao"
	"blog/middleware/hotkey"
	"context"

	"github.com/sirupsen/logrus"
)

type user struct {
	cache *hotkey.HotKeyWithCache
}

func init() {
	userService = newUser()
}

var userService *user

func GetUser() *user {
	return userService
}

func newUser() *user {
	userhotkey, err := hotkey.NewHotkey(&hotkey.Option{
		HotKeyCnt: 50000,
		AutoCache: false,
	})
	if err != nil {
		logrus.Panicf("init hotkey failed:%s", err.Error())
	}
	return &user{
		cache: userhotkey,
	}
}

func (u *user) AutoCreateIfNotExist(ctx context.Context, address string, alias string) (err error) {
	if _, ok := u.cache.Get(address); ok {
		u.cache.Add(address, 1)
		return nil
	}
	userdao := dao.GetUser()
	var user dao.User
	defer func() {
		if err == nil {
			u.cache.Add(address, 1)
		}
	}()
	user, err = userdao.FindUserByAddress(ctx, address)
	if err != nil {
		logrus.Errorf("find user %v failed: %s", address, err.Error())
		return
	}
	if user == (dao.User{}) {
		err = userdao.CreateUser(&dao.User{
			Address: address,
			Alias:   alias,
		})
		if err != nil {
			logrus.Errorf("create user %v failed: %s", address, err.Error())
			return
		}
	}
	return
}
