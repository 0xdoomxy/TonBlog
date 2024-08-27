package utils

import (
	"blog/dao"
	"blog/rpc/kitex_gen/rpc"
	"gorm.io/gorm"
	"time"
)

/*
*
用于包装通过rpc传递的article参数,主要是赋一些默认值
*/
func WrappedArticle(article *rpc.Article) *dao.Article {
	var res = new(dao.Article)
	res.Title = article.Title
	res.Content = article.Content
	res.Creator = article.Creator
	res.Images = article.Images
	res.ID = uint(article.Id)
	res.Tags = article.Tags
	if article.CreateTime > 0 {
		res.CreatedAt = time.Unix(article.CreateTime, 0)
	}
	if article.UpdateTime > 0 {
		res.UpdatedAt = time.Unix(article.UpdateTime, 0)
	}
	if article.DeleteTime > 0 {
		res.DeletedAt = gorm.DeletedAt{Valid: true, Time: time.Unix(article.DeleteTime, 0)}
	}
	return res
}
