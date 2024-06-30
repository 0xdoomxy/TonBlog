package controller

import (
	"blog/service"
	"blog/utils"

	"github.com/gin-gonic/gin"
)

type tag struct {
}

var tagController = &tag{}

func GetTag() *tag {
	return tagController
}

func (t *tag) GetAllTags(ctx *gin.Context) {
	tags, err := service.GetTag().GetTags(ctx)
	if err != nil {
		ctx.JSON(200, utils.NewFailedResponse("获取标签失败"))
		return
	}
	ctx.JSON(200, utils.NewSuccessResponse(tags))
}
