package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

const ContestUserIDKey = "userID"

// getCurrentUser 获取当前登录用户 ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContestUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
