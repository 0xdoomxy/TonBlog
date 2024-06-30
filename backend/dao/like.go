package dao

import (
	"blog/dao/db"
	"context"

	"github.com/spf13/viper"
)

func init() {
	db.GetMysql().AutoMigrate(&Like{})
	likeDao = newLikeDao()
}

type like struct {
	cacheKey string
}

var likeDao *like

func newLikeDao() *like {
	return &like{
		cacheKey: viper.GetString("like.cachekey"),
	}
}

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

func (l *like) IncrementLike(ctx context.Context, like *Like) (err error) {
	cache := db.GetRedis()

	return cache.IncrBy(ctx, l.cacheKey, int64(like.LikeNum)).Err()
}

func (l *like) DecrementLike(ctx context.Context, like *Like) (err error) {
	cache := db.GetRedis()

	return cache.DecrBy(ctx, l.cacheKey, int64(like.LikeNum)).Err()
}

func (l *like) FindLikeById(ctx context.Context, articleid uint) (like Like, err error) {
	err = db.GetMysql().Where("article_id = ?", articleid).First(&like).Error
	return
}
func (l *like) DeleteLike(ctx context.Context, articleid uint) (err error) {
	err = db.GetMysql().Where("article_id = ?", articleid).Delete(&Like{}).Error
	return
}

func (l *like) CreateLike(ctx context.Context, like *Like) (err error) {
	err = db.GetMysql().Create(like).Error
	return
}
