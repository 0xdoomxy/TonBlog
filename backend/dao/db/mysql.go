package db

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", viper.GetString("mysql.username"), viper.GetString("mysql.password"), viper.GetString("mysql.host"), viper.GetInt("mysql.port"), viper.GetString("mysql.database"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panicf("connect to mysql %s failed: %s", dsn, err.Error())
	}

	// Test the database connection
	var t *sql.DB
	t, err = db.DB()
	if err != nil {
		logrus.Panicf("failed to ping database: %s", err.Error())
	}
	err = t.Ping()
	if err != nil {
		logrus.Panicf("failed to ping database: %s", err.Error())
	}

	logrus.Info("connect to mysql success")
}

func GetMysql() *gorm.DB {
	return db
}
