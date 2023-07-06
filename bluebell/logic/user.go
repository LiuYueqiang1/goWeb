package logic

import (
	"bluebell.com/bluebell/dao/mysql"
	"bluebell.com/bluebell/models"
	"bluebell.com/bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParmSignUp) (err error) {
	//1、判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		// 数据库查询出错
		return err
	}
	//2、生成UID
	userID := snowflake.GenID()
	// 构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3、保存进数据库
	return mysql.InsertUser(user)
}
