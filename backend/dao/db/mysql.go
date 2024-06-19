package db

import "gorm.io/gorm"

var mysql *gorm.DB

func GetMysql() *gorm.DB {
	return mysql
}
