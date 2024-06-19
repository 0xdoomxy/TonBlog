package dao

import (
	"blog/dao/db"

	"gorm.io/gorm"
)

func init() {
	db.GetMysql().AutoMigrate(&User{})
}

type user struct {
}

var userDao = &user{}

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

func (u *user) CreateUser(user *User) (err error) {
	err = db.GetMysql().Model(&User{}).Create(user).Error
	if err != nil {
		return
	}
	return
}

func (u *user) FindUserById(id uint) (user User, err error) {
	err = db.GetMysql().Model(&User{}).Where("id = ?", id).First(&user).Error
	return
}

func (u *user) DeleteUser(id uint) (err error) {
	err = db.GetMysql().Model(&User{}).Where("id = ?", id).Delete(&User{}).Error
	return
}

func (u *user) UpdateUser(user *User) (err error) {
	err = db.GetMysql().Model(&User{}).Where("id = ?", user.ID).Updates(user).Error
	return
}
