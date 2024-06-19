package dao

import "blog/dao/db"

func init() {
	db.GetMysql().AutoMigrate(&LikeRelationship{})
}

type likeRelationship struct {
}

var likeRelationshipDao = &likeRelationship{}

func GetLikeRelationship() *likeRelationship {
	return likeRelationshipDao
}

/*
文章关注表
*/
type LikeRelationship struct {
	ArticleID uint `gorm:"not null;uniqueIndex:search"`
	UserID    uint `gorm:"not null;uniqueIndex:search"`
}

func (l *likeRelationship) CreateLikeRelationship(likeRelationship *LikeRelationship) (err error) {
	err = db.GetMysql().Model(&LikeRelationship{}).Create(likeRelationship).Error
	if err != nil {
		return
	}
	return
}

func (l *likeRelationship) FindLikeRelationshipById(articleID, userID uint) (likeRelationship LikeRelationship, err error) {
	err = db.GetMysql().Model(&LikeRelationship{}).Where("article_id = ? and user_id = ?", articleID, userID).First(&likeRelationship).Error
	return
}

func (l *likeRelationship) DeleteLikeRelationship(articleID, userID uint) (err error) {
	err = db.GetMysql().Model(&LikeRelationship{}).Where("article_id = ? and user_id = ?", articleID, userID).Delete(&LikeRelationship{}).Error
	return
}

func (l *likeRelationship) UpdateLikeRelationship(likeRelationship *LikeRelationship) (err error) {
	err = db.GetMysql().Model(&LikeRelationship{}).Where("article_id = ? and user_id = ?", likeRelationship.ArticleID, likeRelationship.UserID).Updates(likeRelationship).Error
	return
}

func (l *likeRelationship) FindLikeRelationshipByArticleID(articleID uint) (likeRelationship []LikeRelationship, err error) {
	err = db.GetMysql().Model(&LikeRelationship{}).Where("article_id = ?", articleID).Find(&likeRelationship).Error
	return
}
