package dao

import "blog/dao/db"

func init() {
	db.GetMysql().AutoMigrate(&Comment{})
}

type comment struct {
	_ [0]func()
}

var commentDao = &comment{}

func GetComment() *comment {
	return commentDao
}

/*
*
	评论表
*
*/
type Comment struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	CreateAt  uint64 `gorm:"autoCreateTime:milli"`
	SubID     uint   `gorm:"not null;index:search"`
	Content   string `gorm:"type:varchar(255);not null"`
	ArticleID uint   `gorm:"not null;index:search"`
	Creator   string `gorm:"varchar(64) not null"`
}

func (c *comment) CreateComment(comment *Comment) (err error) {
	err = db.GetMysql().Model(&Comment{}).Create(comment).Error
	if err != nil {
		return
	}
	return
}

func (c *comment) FindCommentById(id uint) (comment Comment, err error) {
	err = db.GetMysql().Model(&Comment{}).Where("id = ?", id).First(&comment).Error
	return
}

func (c *comment) DeleteComment(id uint) (err error) {
	err = db.GetMysql().Model(&Comment{}).Where("id = ?", id).Delete(&Comment{}).Error
	return
}

func (c *comment) FindCommentByArticleId(articleId uint) (comments []Comment, err error) {
	err = db.GetMysql().Model(&Comment{}).Where("article_id = ?", articleId).Find(&comments).Error
	return
}

func (c *comment) FindCommentBySubId(subId uint) (comments []Comment, err error) {
	err = db.GetMysql().Model(&Comment{}).Where("sub_id = ?", subId).Find(&comments).Error
	return
}

func (c *comment) UpdateComment(comment *Comment) (err error) {
	err = db.GetMysql().Model(&Comment{}).Where("id = ?", comment.ID).Updates(comment).Error
	return
}

func (c *comment) DeleteCommentByArticleId(articleId uint) (err error) {
	err = db.GetMysql().Model(&Comment{}).Where("article_id = ?", articleId).Delete(&Comment{}).Error
	return
}
