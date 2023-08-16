package routes

import (
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"

	"bluebell.com/bluebell/controllers"
	"bluebell.com/bluebell/middlewares"

	"bluebell.com/bluebell/logger"
	"github.com/gin-gonic/gin"

	_ "bluebell.com/bluebell/docs" // 千万不要忘了导入把你上一步生成的docs
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置程发布者模式
	}

	r := gin.New()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.LoadHTMLFiles("./templates/index.html")
	r.Static("static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// 两秒钟限制
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
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
		v1.GET("/posts2", controllers.GetPostListHandler2)
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

// http://127.0.0.1:8081/swagger/index.html
