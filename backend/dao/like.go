package dao

import "blog/dao/db"

func init() {
	db.GetMysql().AutoMigrate(&Like{})
}

type like struct {
}

var likeDao = &like{}

func GetLike() *like {
	return likeDao
}

/*
*

	文章关注总览表

*
*/
type Like struct {
	ArticleID uint `gorm:"not null;uniqueIndex:search"`
	LikeNum   uint `gorm:"not null"`
}
