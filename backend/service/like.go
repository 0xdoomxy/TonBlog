package service

import (
	"blog/dao"
	"blog/model"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type like struct {
}

var likeService = &like{}

func GetLike() *like {
	return likeService
}

func (l *like) SetAsLike(ctx context.Context, address string, articleid uint) (err error) {
	articledao := dao.GetArticle()
	_, err = articledao.FindArticleById(ctx, articleid)
	if err != nil {
		return
	}
	likeDAO := dao.GetLike()
	like_relationshipDAO := dao.GetLikeRelationship()
	var ok bool
	ok, err = like_relationshipDAO.FindLikeRelationshipByArticleIDAndUserid(ctx, &model.LikeRelationship{ArticleID: articleid, Address: address})
	if ok || err != nil {
		if err != nil {
			logrus.Errorf("find like relationship by articleid and userid failed: %v", err)
		} else {
			err = fmt.Errorf("repeat cancel like")
		}
		return
	}
	like := &model.Like{
		ArticleID: articleid,
		LikeNum:   1,
	}
	err = likeDAO.IncrementLike(ctx, like)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			likeDAO.DecrementLike(ctx, like)
		}
	}()
	like_relationship := &model.LikeRelationship{
		ArticleID: articleid,
		Address:   address,
	}
	err = like_relationshipDAO.CreateLikeRelationship(ctx, like_relationship)
	return
}

// 是否重复点赞或取消点赞
func (l *like) CancelLike(ctx context.Context, address string, articleid uint) (err error) {
	articledao := dao.GetArticle()
	_, err = articledao.FindArticleById(ctx, articleid)
	if err != nil {
		return
	}
	likeDAO := dao.GetLike()
	like_relationshipDAO := dao.GetLikeRelationship()
	var ok bool
	ok, err = like_relationshipDAO.FindLikeRelationshipByArticleIDAndUserid(ctx, &model.LikeRelationship{ArticleID: articleid, Address: address})
	if ok || err != nil {
		if err != nil {
			logrus.Errorf("find like relationship by articleid and userid failed: %v", err)
		} else {
			err = fmt.Errorf("repeat cancel like")
		}
		return
	}
	like := &model.Like{
		ArticleID: articleid,
		LikeNum:   1,
	}
	err = likeDAO.DecrementLike(ctx, like)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			likeDAO.IncrementLike(ctx, like)
		}
	}()
	like_relationship := &model.LikeRelationship{
		ArticleID: articleid,
		Address:   address,
	}
	err = like_relationshipDAO.DeleteLikeRelationship(ctx, like_relationship)
	return
}

func (l *like) FindIsExist(ctx context.Context, articleid uint, address string) (exist bool, err error) {
	like_relationshipDap := dao.GetLikeRelationship()
	exist, err = like_relationshipDap.FindLikeRelationshipByArticleIDAndUserid(ctx, &model.LikeRelationship{ArticleID: articleid, Address: address})
	return
}
