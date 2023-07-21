package controllers

import (
	"bluebell.com/bluebell/logic"
	"bluebell.com/bluebell/models"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	//1、获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 由于始终拿不到 token ，所以我们取消了登录这一项

	// 从 c 取到当前发送请求的用户的 ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	//2、创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3、返回响应
	ResponseSuccess(c, nil)
}

/*
"title":"study go make me happy",
"content":"My ei must be accept",
"community_id":1

   "username":"七米",
   "password":"123"
*/
