package controller

import (
	"blog/service"
	"blog/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type comment struct {
}

var commentController *comment

func GetComment() *comment {
	return commentController
}

func (c *comment) CreateComment(ctx *gin.Context) {
	creator, ok := ctx.Get("publickey")
	if !ok {
		ctx.JSON(401, utils.NewFailedResponse("未登录"))
		return
	}
	var err error
	var comment struct {
		ArticleID uint   `json:"articleid"`
		Content   string `json:"content"`
		TopID     int    `json:"topid"`
	}
	err = ctx.BindJSON(&comment)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	err = service.GetComment().CreateComment(ctx, uint(comment.ArticleID), creator.(string), comment.Content, comment.TopID)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("评论失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(nil))
}

func (c *comment) FindCommentByArticle(ctx *gin.Context) {
	articleidStr := ctx.Query("articleid")
	articleid, err := strconv.ParseUint(articleidStr, 10, 64)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	comments, err := service.GetComment().FindCommentByArticle(ctx, uint(articleid))
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("查询失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(comments))
}

func (c *comment) DeleteComment(ctx *gin.Context) {
	creator := ctx.GetString("publickey")
	articleidStr := ctx.Query("articleid")
	articleid, err := strconv.ParseUint(articleidStr, 10, 64)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("参数错误"))
		return
	}
	err = service.GetComment().DeleteComment(ctx, uint(articleid), uint(id), creator)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("删除失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(nil))
}
