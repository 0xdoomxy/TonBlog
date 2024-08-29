package service

import (
	"blog/dao"
	"blog/model"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type tag struct {
}

var tagService = &tag{}

func GetTag() *tag {
	return tagService
}

type TagsView dao.Tags

func (t *tag) GetTags(ctx context.Context) (view TagsView, err error) {
	var tags dao.Tags
	tags, err = dao.GetTag().FindAllTags(ctx)
	return TagsView(tags), err
}
func (t *tag) CreateTag(ctx context.Context, tag *model.Tag) (err error) {
	err = dao.GetTag().CreateTag(ctx, tag)
	return
}

func (t *tag) CreateTagRelationship(ctx context.Context, tagRelationship *model.TagRelationship) (err error) {
	err = dao.GetTagRelationship().CreateTagRelationship(ctx, tagRelationship)
	return
}
func (t *tag) DeleteTagRelationship(ctx context.Context, tagRelationship *model.TagRelationship) (err error) {
	err = dao.GetTagRelationship().DeleteTagRelationship(ctx, tagRelationship)
	return
}

/*
*

	根据标签名字来增加文章数量

*
*/
func (t *tag) IncrementArticleNumByName(ctx context.Context, name string, num uint) (err error) {
	tagDao := dao.GetTag()
	err = tagDao.FindAndIncrementTagNumByName(ctx, name, num)
	if err != nil {
		logrus.Errorf("find and increment tag num by name (%s) failed: %s", name, err.Error())
		return
	}
	return
}

func (t *tag) IncrementArticleNumByNames(ctx context.Context, names []string, num uint) (err error) {
	for _, name := range names {
		err = t.IncrementArticleNumByName(ctx, name, num)
		if err != nil {
			return
		}
	}
	return
}
func (t *tag) GetArticleTotalByName(ctx context.Context, name string) (total uint, err error) {
	var rsv model.Tag
	rsv, err = dao.GetTag().FindTag(ctx, name)
	if err != nil {
		logrus.Errorf("find tag by name (%s) failed: %s", name, err.Error())
	}
	return rsv.ArticleNum, err
}

func (t *tag) FindArticlesByTagName(ctx context.Context, name string, page int, pagesize int) (view *ArticleViewByPage, err error) {
	var articles dao.TagRelationships
	view = new(ArticleViewByPage)
	articles, err = dao.GetTagRelationship().FindTagRelationshipByName(ctx, name, page, pagesize)
	if err != nil {
		logrus.Errorf("find articles by tag name (%s) failed: %s", name, err.Error())
		return
	}
	tagService := GetTag()
	view.Total, err = tagService.GetArticleTotalByName(ctx, name)
	if err != nil {
		logrus.Errorf("get article total by name (%s) failed: %s", name, err.Error())
		return
	}
	articleService := GetArticle()
	wg := sync.WaitGroup{}
	wg.Add(len(articles))
	for _, article := range articles {
		go func(id uint) {
			defer wg.Done()
			if err != nil {
				return
			}
			var a *ArticleView
			a, err = articleService.FindArticlePatical(ctx, id)
			if err != nil {
				logrus.Errorf("find article by id (%d) failed: %s", id, err.Error())
				return
			}
			view.Articles = append(view.Articles, a)
		}(article.ArticleId)
	}
	wg.Wait()
	return
}
