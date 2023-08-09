package routes

import (
	"net/http"

	"bluebell.com/bluebell/controllers"
	"bluebell.com/bluebell/middlewares"

	"bluebell.com/bluebell/logger"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置程发布者模式
	}

	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", controllers.SignUpHandler)

	// 登录
	v1.POST("/login", controllers.LoginHandler)
	v1.GET("/community", controllers.CommunityHandler)
	//v1.GET("/community/:id", controllers.CommunityDetailHandler)
	v1.GET("/community/:id", controllers.CommunityDetailHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JTW认证中间件
	{
		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/posts", controllers.GetPostListHandler)
		// http://127.0.0.1:8081/api/v1/posts/?size=1&page=3

		//投票
		v1.POST("/vote", controllers.PostVoteController)
	}
	//
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是登录的用户，判断请求头中是否有 有效的JWT？
		c.String(http.StatusOK, "ping")
	})
	//r.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "ok")
	//})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
