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
	// **********
	// 将网页中拿到的用户名密码反序列化到 models.ParmSignUp 结构体中
	// 再将 models.ParmSignUp 中的值赋给定义的 models.User 的结构体中
	// **********
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3、保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParmLogin) (err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return err
	}
	return
}
