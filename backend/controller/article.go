package controller

import (
	"blog/dao"
	"blog/service"
	"blog/utils"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type article struct {
}

var articleController = &article{}

func GetArticle() *article {
	return articleController
}

func (a *article) PublishArticle(ctx *gin.Context) {
	// userid := context.GetInt64("userid")
	var publickey, ok = ctx.Get("publickey")
	if !ok {
		ctx.JSON(200, utils.NewFailedResponse("未登录"))
		return
	}
	var err error
	var form *multipart.Form
	form, err = ctx.MultipartForm()
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	var article = &dao.Article{}
	//添加userid 到articleMap中,必须确保articleMap中没有userid字段
	article.Creator = publickey.(string)
	article.Title = form.Value["title"][0]
	article.Content = form.Value["content"][0]
	article.Images = form.Value["images"][0]
	article.Tags = form.Value["tags"][0]
	err = service.GetArticle().PublishArticle(ctx, article)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("发布失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(nil))
}

func (a *article) UploadImage(context *gin.Context) {
	//TODO this is a dangerous code, you should check the user's permission
	userid := 1
	if userid <= 0 {
		context.JSON(200, utils.NewFailedResponse("未登录"))
		return
	}
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(200, utils.NewFailedResponse("上传失败"))
	}
	if !utils.IsImage(file.Filename) {
		context.JSON(200, utils.NewFailedResponse("不是图片"))
		return
	}
	imgtype := filepath.Ext(file.Filename)
	var f multipart.File
	f, err = file.Open()
	if err != nil {
		context.JSON(200, utils.NewFailedResponse("上传失败"))
	}
	fileName := fmt.Sprintf("%d_%d%s", userid, time.Now().Unix(), imgtype)
	err = service.GetArticle().UploadImage(fileName, f)
	if err != nil {
		context.JSON(200, utils.NewFailedResponse("上传失败"))
	}
	context.JSON(200, utils.NewSuccessResponse(fileName))
}

func (a *article) DownloadImage(context *gin.Context) {
	filename := context.Query("filename")
	if filename == "" {
		context.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	res, err := service.GetArticle().DownloadImage(filename)
	if err != nil {
		context.JSON(200, utils.NewFailedResponse("下载失败"))
		return
	}
	context.Writer.Write(res)
}

func (a *article) FindArticle(ctx *gin.Context) {
	var err error
	articleid, err := strconv.ParseUint(ctx.Query("id"), 10, 64)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	var view *service.ArticleView
	view, err = service.GetArticle().FindArticle(ctx, uint(articleid))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(view))
}

/*
*
查询讨论数最多的文章
*
*/
func (a *article) FindArticleByMaxAccessNum(ctx *gin.Context) {
	var err error
	var view *service.ArticleViewByPage
	var page int
	var pageSize int
	page, err = strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	pageSize, err = strconv.Atoi(ctx.Query("pagesize"))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	view, err = service.GetArticle().FindArticleByAccessNum(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(view))
}

/**
查询最新的文章
**/

func (a *article) FindArticleByCreateTime(ctx *gin.Context) {
	var page int
	var pageSize int
	var err error

	page, err = strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	pageSize, err = strconv.Atoi(ctx.Query("pagesize"))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	var view *service.ArticleViewByPage
	view, err = service.GetArticle().FindArticlePaticalByCreateTime(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(view))
}

/*
*
根据关键字搜索文章

*
*/
func (a *article) SearchArticleByPage(ctx *gin.Context) {
	var err error
	var view *service.ArticleViewByPage
	var page int
	var pageSize int
	page, err = strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	pageSize, err = strconv.Atoi(ctx.Query("pagesize"))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	view, err = service.GetArticle().SearchArticleByPage(ctx, ctx.Query("keyword"), page, pageSize)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(view))
}
