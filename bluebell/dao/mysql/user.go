package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	"bluebell.com/bluebell/models"
)

const serect = "liwenzhou.com"

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

// 判断用户存不存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// Login
func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	// var ErrNoRows = errors.New("sql: no rows in result set")
	if err == sql.ErrNoRows {
		// 打印用户不存在 已经定义为全局变量
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		// 如果密码不相等，打印用户名或密码错误 已经定义为全局变量
		return ErrorInvalidPassword
	}
	return
}

// 对密码进行加密
func encryPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(serect))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(serect))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
