package controller

import (
	"blog/service"
	"blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type like struct {
}

var likeController = &like{}

func GetLike() *like {
	return likeController
}

func (l *like) SetAsLike(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Query("articleid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("参数出错"))
		return
	}
	userId, err := strconv.Atoi(ctx.Query("userid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("参数出错"))
		return
	}
	err = service.GetLike().SetAsLike(ctx, uint(userId), uint(articleId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("点赞失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(nil))
}

func (l *like) CancelLike(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Query("articleid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("参数出错"))
		return
	}
	userId, err := strconv.Atoi(ctx.Query("userid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("参数出错"))
		return
	}
	err = service.GetLike().CancelLike(ctx, uint(userId), uint(articleId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("取消点赞失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(nil))
}

func (l *like) IsExist(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Query("articleid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("参数出错"))
		return
	}
	userId, err := strconv.Atoi(ctx.Query("userid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("参数出错"))
		return
	}
	var exist bool
	exist, err = service.GetLike().FindIsExist(ctx, uint(articleId), uint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewFailedResponse("取消点赞失败"))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(struct {
		Exist bool `json:"exist"`
	}{
		Exist: exist,
	}))
}
