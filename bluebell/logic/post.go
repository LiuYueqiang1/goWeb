package logic

import (
	"bluebell.com/bluebell/dao/mysql"
	"bluebell.com/bluebell/models"
	"bluebell.com/bluebell/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	//err = mysql.CreatePost(p)
	//if err != nil {
	//	return err
	//}
	//err = redis.CreatePost(p.ID, p.CommunityID)
	return mysql.CreatePost(p)
	// 3. 返回
}
