package controllers

import (
	"strconv"

	"bluebell.com/bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//---跟社区相关---

func CommunityHandler(c *gin.Context) {
	//查询到所有的社区以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//1、获取社区id
	idStr := c.Param("id") //获取URL参数
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//查询到所有的社区以列表的形式返回
	//2、根据id获取社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
