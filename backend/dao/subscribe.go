package dao

import (
	"blog/dao/db"
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetSubscribeDao() *subscribe {
	return subscribeDao
}

type subscribe struct {
	_ [0]func()
}

var subscribeDao = new(subscribe)

func (s *subscribe) CreateSubscribe(ctx context.Context, sub *Subscribe) (id uint, err error) {
	conn := db.GetMysql()
	err = conn.Model(&Subscribe{}).WithContext(ctx).Create(sub).Error
	if err != nil {
		logrus.Error("create subscribe err:%v", err)
		return
	}
	id = sub.ID
	//cache := db.GetRedis()
	return
}

func (s *subscribe) CancelSubscribe(ctx context.Context, id uint) (err error) {
	conn := db.GetMysql()
	err = conn.Model(&Subscribe{}).WithContext(ctx).Where("id = ?", id).Delete(&Comment{}).Error
	if err != nil {
		logrus.Error("cancel subscribe err:%v", err)
	}
	return err
}

func (s *subscribe) QuerySubscribeByType(ctx context.Context, typ SubscribeType) (res *[]Subscribe, err error) {
	conn := db.GetMysql()
	err = conn.Model(&Subscribe{}).WithContext(ctx).Where("type = ?", typ).Find(&res).Error
	if err != nil {
		logrus.Error("query  subscribe by type err:%v", err)
	}
	return
}

func (s *subscribe) QuerySubscribeByCreatorAndType(ctx context.Context, creator uint, typ SubscribeType) (res *[]Subscribe, err error) {
	conn := db.GetMysql()
	err = conn.Model(&Subscribe{}).WithContext(ctx).Where("creator = ? and type = ?", creator, typ).Find(&res).Error
	if err != nil {
		logrus.Error("query  subscribe by creator and subscribe type  err:%v", err)
	}
	return
}

type Subscribes []*Subscribe

func (subscribes *Subscribes) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, subscribes)
}
func (subscribes *Subscribes) MarshalBinary() ([]byte, error) {
	return json.Marshal(subscribes)
}

type SubscribeType int

const (
	SubscribeTag SubscribeType = iota
	SubscribeAuthor
)

func init() {
	db.GetMysql().AutoMigrate(&Subscribe{})
}

type Subscribe struct {
	gorm.Model
	Creator int64         `gorm:"not null;index:idx_find"`
	Type    SubscribeType `gorm:"index:idx_type;not null;index:idx_find"`
	//split up the comma symbol
	Object string `gorm:"type:varchar(50);not null"`
}

func (s *Subscribe) TableName() string {
	return "subscribe"
}

func (s *Subscribe) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
func (s *Subscribe) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
