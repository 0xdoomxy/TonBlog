package service

import (
	"blog/dao"
	"context"

	"github.com/sirupsen/logrus"
)

type tag struct {
}

var tagService = &tag{}

func GetTag() *tag {
	return tagService
}

type TagsView []*dao.Tag

func (t *tag) GetTags(ctx context.Context) (view TagsView, err error) {
	var tags []*dao.Tag
	tags, err = dao.GetTag().FindAllTags(ctx)
	return TagsView(tags), err
}
func (t *tag) CreateTag(ctx context.Context, tag *dao.Tag) (err error) {
	err = dao.GetTag().CreateTag(ctx, tag)
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
