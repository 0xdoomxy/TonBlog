package dao

import (
	"blog/dao/db"
	"context"

	"gorm.io/gorm"
)

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

func (t *tag) CreateTag(ctx context.Context, tag *Tag) (err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Create(tag).Error
	if err != nil {
		return
	}
	return
}

func (t *tag) FindTagById(ctx context.Context, id uint) (tag Tag, err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("id = ?", id).First(&tag).Error
	return
}

func (t *tag) DeleteTag(ctx context.Context, id uint) (err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("id = ?", id).Delete(&Tag{}).Error
	return
}

func (t *tag) UpdateTag(ctx context.Context, tag *Tag) (err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("id = ?", tag.ID).Updates(tag).Error
	return
}

func (t *tag) FindTagByName(ctx context.Context, name string) (tag Tag, err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("name = ?", name).First(&tag).Error
	return
}

func (t *tag) FindAndIncrementTagNumById(ctx context.Context, id uint, num uint) (err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Where("id = ?", id).Update("article_num", gorm.Expr("article_num + ?", num)).Error
	return
}
func (t *tag) FindAndIncrementTagNumByName(ctx context.Context, name string, num uint) (err error) {
	//乐观
	result := db.GetMysql().WithContext(ctx).Model(&Tag{}).Where("name = ?", name).Update("article_num", gorm.Expr("article_num + ?", num))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return db.GetMysql().Model(&Tag{}).Create(&Tag{Name: name, ArticleNum: num}).Error
	}
	return
}
func (t *tag) FindAllTags(ctx context.Context) (tags []*Tag, err error) {
	err = db.GetMysql().Model(&Tag{}).WithContext(ctx).Find(&tags).Error
	return
}
