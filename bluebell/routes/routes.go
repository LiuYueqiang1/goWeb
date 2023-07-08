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

	// 注册
	r.POST("/signup", controllers.SignUpHandler)

	// 登录
	r.POST("/login", controllers.LoginHandler)

	//
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是登录的用户，判断请求头中是否有 有效的JWT？
		c.String(http.StatusOK, "ping")
	})
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
