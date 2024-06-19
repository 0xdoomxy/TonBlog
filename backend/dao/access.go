package dao

import "blog/dao/db"

func init() {
	db.GetMysql().AutoMigrate(&Access{})
}

type access struct {
}

var accessDao = &access{}

/**
	访问表
**/
type Access struct {
	ArticleID uint `gorm:"primaryKey"`
	AccessNum uint `gorm:"not null"`
}

func GetAccess() *access {
	return accessDao
}

func (a *access) CreateAccess(access *Access) (err error) {
	err = db.GetMysql().Model(&Access{}).Create(access).Error
	if err != nil {
		return
	}
	return
}

func (a *access) FindAccessById(id uint) (access Access, err error) {
	err = db.GetMysql().Model(&Access{}).Where("article_id = ?", id).First(&access).Error
	return
}

func (a *access) DeleteAccess(id uint) (err error) {
	err = db.GetMysql().Model(&Access{}).Where("article_id = ?", id).Delete(&Access{}).Error
	return
}
