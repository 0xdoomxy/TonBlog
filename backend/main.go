package main

import (
	"blog/controller"
	"blog/cron"
	"blog/middleware/cors"
	"blog/middleware/jwt"
	"blog/middleware/metrics"
	"blog/middleware/whitepaper"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func StartCronTask() {
	manager := cron.NewCronManager()
	manager.EquipmentTask(cron.NewAccessConsumerCron(), cron.NewLikeConsumerCron())
	go manager.Run()
}

func main() {
	//starting the cron task
	StartCronTask()
	//registe  gin router
	engine := gin.Default()
	engine.Use(cors.CORS())
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	engine.Use(gin.Recovery())
	//绑定流量监控
	metrics.BindMetrics(engine)
	//绑定业务路由
	bindArticleRoutes(engine)
	bindLikeRoutes(engine)
	bindCommentRoutes(engine)
	bindUserRoutes(engine)
	bindTagRoutes(engine)
	bindAirportRoutes(engine)
	engine.Run(":8080")
}
func bindAirportRoutes(engine *gin.Engine) {
	route := engine.Group("/airport")
	route.GET("/findrunning", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetAirport().FindRunningAirport(ctx)
	})
	route.GET("/findfinish", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetAirport().FindFinishAirport(ctx)
	})
	route.POST("/create", jwt.NewVerifyMiddleware(), whitepaper.WhitepaperMiddleware(), func(ctx *gin.Context) {
		controller.GetAirport().CreateAirport(ctx)
	})
	route.GET("/delete", jwt.NewVerifyMiddleware(), whitepaper.WhitepaperMiddleware(), func(ctx *gin.Context) {
		controller.GetAirport().DeleteAirport(ctx)
	})
	route.GET("/update", jwt.NewVerifyMiddleware(), func(context *gin.Context) {
		controller.GetAirport().UpdateAirport(context)
	})
}

func bindArticleRoutes(engine *gin.Engine) {
	route := engine.Group("/article")

	route.POST("/image/upload", jwt.NewVerifyMiddleware(), func(context *gin.Context) {
		controller.GetArticle().UploadImage(context)
	})
	route.GET("/image/download", func(ctx *gin.Context) {
		controller.GetArticle().DownloadImage(ctx)
	})
	route.POST("/publish", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetArticle().PublishArticle(ctx)
	})
	route.GET("/findbymaxaccess", func(ctx *gin.Context) {
		controller.GetArticle().FindArticleByMaxAccessNum(ctx)
	})
	route.GET("/findbycreatetime", func(ctx *gin.Context) {
		controller.GetArticle().FindArticleByCreateTime(ctx)
	})
	route.GET("/find", func(ctx *gin.Context) {
		controller.GetArticle().FindArticle(ctx)
	})
	route.GET("/search", func(ctx *gin.Context) {
		controller.GetArticle().SearchArticleByPage(ctx)
	})
	route.GET("/delete", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetArticle().DeleteArticle(ctx)
	})
	route.GET("/update", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetArticle().UpdateArticle(ctx)
	})
}

func bindLikeRoutes(engine *gin.Engine) {
	router := engine.Group("/like")
	router.GET("/confirm", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetLike().SetAsLike(ctx)
	})
	router.GET("/cancel", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetLike().CancelLike(ctx)
	})
	router.GET("/exist", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetLike().IsExist(ctx)
	})
}
func bindCommentRoutes(engine *gin.Engine) {
	router := engine.Group("/comment")
	router.POST("/create", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetComment().CreateComment(ctx)
	})
	router.GET("/find", func(ctx *gin.Context) {
		controller.GetComment().FindCommentByArticle(ctx)
	})
	router.GET("/delete", jwt.NewVerifyMiddleware(), func(ctx *gin.Context) {
		controller.GetComment().DeleteComment(ctx)
	})

}

// TODO test the proof is true
func bindUserRoutes(engine *gin.Engine) {
	router := engine.Group("/user")
	router.POST("/login", func(ctx *gin.Context) {
		controller.GetUser().LoginHandler(ctx)
	})
}

func bindTagRoutes(engine *gin.Engine) {
	router := engine.Group("/tag")
	router.GET("/findall", func(ctx *gin.Context) {
		controller.GetTag().GetAllTags(ctx)
	})
	router.GET("/findArticle", func(ctx *gin.Context) {
		controller.GetTag().GetArticleByTag(ctx)
	})
}
