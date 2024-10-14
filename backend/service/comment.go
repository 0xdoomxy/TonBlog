package service

import (
	"blog/dao"
	"blog/model"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type comment struct {
}

var commentService = &comment{}

func GetComment() *comment {
	return commentService
}

func (c *comment) CreateComment(ctx context.Context, articleid uint, creator string, content string, topid int) (err error) {
	articledao := dao.GetArticle()
	_, err = articledao.FindArticleById(ctx, articleid)
	if err != nil {
		return
	}
	comment := model.Comment{
		ArticleID: articleid,
		Creator:   creator,
		Content:   content,
		TopID:     uint(topid),
		CreateAt:  time.Now(),
	}
	err = dao.GetComment().CreateComment(ctx, &comment)
	return
}
func (c *comment) FindCommentByArticle(ctx context.Context, articleid uint) (comments []*model.Comment, err error) {
	comments, err = dao.GetComment().FindCommentByArticleid(ctx, articleid)
	if err != nil {
		logrus.Error("find comment by article failed:", err)
	}
	return
}
func (c *comment) DeleteComment(ctx context.Context, articleid uint, id uint, creator string) (err error) {
	commentdao := dao.GetComment()
	articledao := dao.GetArticle()
	var ok bool
	var article model.Article
	article, err = articledao.FindArticlePaticalById(ctx, articleid)
	if err != nil {
		logrus.Errorf("find article failed:%v", err)
		return
	}
	if article.Creator == creator {
		ok = true
	}
	if !ok {
		ok, err = commentdao.FindCommentCreateBy(ctx, id, creator)
		if !ok {
			err = fmt.Errorf("cant delete others comment")
			return
		}
	}
	if ok {
		err = dao.GetComment().DeleteComment(ctx, articleid, id)
		if err != nil {
			logrus.Error("delete comment failed:", err)
		}
	}
	return
}
