package routes

import (
	"net/http"

	"bluebell.com/bluebell/controllers"

	"bluebell.com/bluebell/logger"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册
	r.POST("/signup", controllers.SignUpHandler)

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
