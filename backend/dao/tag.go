package dao

import "blog/dao/db"

func init() {
	db.GetMysql().AutoMigrate(&Tag{})
}

type tag struct {
}

var tagDao = &tag{}

func GetTag() *tag {
	return tagDao
}

/*
*

	标签表

*
*/
type Tag struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(255);not null;index:search"`
	ArticleNum uint   `gorm:"not null"`
}

func (t *tag) CreateTag(tag *Tag) (err error) {
	err = db.GetMysql().Model(&Tag{}).Create(tag).Error
	if err != nil {
		return
	}
	return
}

func (t *tag) FindTagById(id uint) (tag Tag, err error) {
	err = db.GetMysql().Model(&Tag{}).Where("id = ?", id).First(&tag).Error
	return
}

func (t *tag) DeleteTag(id uint) (err error) {
	err = db.GetMysql().Model(&Tag{}).Where("id = ?", id).Delete(&Tag{}).Error
	return
}

func (t *tag) UpdateTag(tag *Tag) (err error) {
	err = db.GetMysql().Model(&Tag{}).Where("id = ?", tag.ID).Updates(tag).Error
	return
}

func (t *tag) FindTagByName(name string) (tag Tag, err error) {
	err = db.GetMysql().Model(&Tag{}).Where("name = ?", name).First(&tag).Error
	return
}
