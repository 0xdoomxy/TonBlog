package dao

import (
	"blog/dao/db"
	"bytes"
	"io"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var articleContentBucketName = "0xdoomxy-blog"

func init() {
	db.GetMysql().AutoMigrate(&Article{})
}

type article struct {
}

var articleDao = &article{}

func GetArticle() *article {
	return articleDao
}

/*
*
文章表
*
*/
type Article struct {
	gorm.Model
	Title   string `gorm:"type:varchar(255);not null"`
	Tags    string `gorm:"type:varchar(255)"`
	Creator uint   `gorm:"not null"`
	Summary string `gorm:"type:varchar(255)"`
	Content []byte `gorm:"-"`
}

func (a *article) CreateArticle(article *Article) (err error) {
	err = db.GetBucket(articleContentBucketName).PutObject(strconv.Itoa(int(article.ID)), bytes.NewReader(article.Content))
	defer func() {
		if err != nil {
			db.GetBucket(articleContentBucketName).DeleteObject(strconv.Itoa(int(article.ID)))
		}
	}()
	err = db.GetMysql().Model(&Article{}).Create(article).Error
	if err != nil {
		return
	}
	return
}

func (a *article) UpdateArticle(article *Article) (err error) {
	err = db.GetBucket(articleContentBucketName).PutObject(strconv.Itoa(int(article.ID)), bytes.NewReader(article.Content))
	if err != nil {
		return
	}
	err = db.GetMysql().Model(&Article{}).Where("id = ?", article.ID).Updates(article).Error
	return
}
func (a *article) DeleteArticle(id uint) (err error) {
	err = db.GetMysql().Model(&Article{}).Where("id = ?", id).Delete(&Article{}).Error
	if err != nil {
		return
	}
	err = db.GetBucket(articleContentBucketName).DeleteObject(strconv.Itoa(int(id)))
	return
}
func (a *article) FindArticlePaticalById(id uint) (article Article, err error) {
	err = db.GetMysql().Model(&Article{}).Where("id = ?", id).First(&article).Error
	return
}
func (a *article) FindArticleById(id uint) (article Article, err error) {
	err = db.GetMysql().Model(&Article{}).Where("id = ?", id).First(&article).Error
	if err != nil {
		return
	}
	var reader io.ReadCloser
	reader, err = db.GetBucket(articleContentBucketName).GetObject(strconv.Itoa(int(article.ID)))
	if err != nil {
		logrus.Error("获取文章内容失败", err.Error())
		return
	}
	defer reader.Close()
	article.Content, err = io.ReadAll(reader)
	if err != nil {
		logrus.Error("读取文章内容失败", err.Error())
		return
	}
	return
}
