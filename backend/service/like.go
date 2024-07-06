package service

import (
	"blog/dao"
	"context"
)

type like struct {
}

var likeService = &like{}

func GetLike() *like {
	return likeService
}

// TODO before set as like ,should make sure the userid is valid
func (l *like) SetAsLike(ctx context.Context, publickey string, articleid uint) (err error) {
	likeDAO := dao.GetLike()
	like_relationshipDAO := dao.GetLikeRelationship()

	like := &dao.Like{
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
	like_relationship := &dao.LikeRelationship{
		ArticleID: articleid,
		PublicKey: publickey,
	}
	err = like_relationshipDAO.CreateLikeRelationship(ctx, like_relationship)
	return
}

// TODO before set as like ,should make sure the userid is valid
func (l *like) CancelLike(ctx context.Context, publickey string, articleid uint) (err error) {
	likeDAO := dao.GetLike()
	like_relationshipDAO := dao.GetLikeRelationship()
	like := &dao.Like{
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
	like_relationship := &dao.LikeRelationship{
		ArticleID: articleid,
		PublicKey: publickey,
	}
	err = like_relationshipDAO.DeleteLikeRelationship(ctx, like_relationship)
	return
}

func (l *like) FindIsExist(ctx context.Context, articleid uint, publickey string) (exist bool, err error) {
	like_relationshipDap := dao.GetLikeRelationship()
	exist, err = like_relationshipDap.FindLikeRelationshipByArticleIDAndUserid(ctx, &dao.LikeRelationship{ArticleID: articleid, PublicKey: publickey})
	return
}
