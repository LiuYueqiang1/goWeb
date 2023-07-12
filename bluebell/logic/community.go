package logic

import (
	"bluebell.com/bluebell/dao/mysql"
	"bluebell.com/bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {

	// 查询数据库 查找到所有的community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
