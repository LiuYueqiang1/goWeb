package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"bluebell.com/bluebell/models"

	"bluebell.com/bluebell/logic"
	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1、获取参数和参数校验
	p := new(models.ParmSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数、有误直接返回响应
		// 记录日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		// 获得应用程序的错误消息的所有信息
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			//"msg": "请求参数有误",
			//查看哪里有误
			//"msg": err.Error(),
			// 翻译
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}

	//` `中使用 binding:"required" 替换以下功能
	// 手动对请求参数进行详细的业务规则校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}

	fmt.Println(p)
	//2、业务处理
	if err := logic.SignUp(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	//3、返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func LoginHandler(c *gin.Context) {
	//1、获取参数和参数校验
	p := new(models.ParmLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数、有误直接返回响应
		// 记录日志
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	//2、业务处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username:", p.Username), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}
	//3、返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
