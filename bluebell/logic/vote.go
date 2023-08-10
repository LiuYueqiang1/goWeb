package logic

import (
	"bluebell.com/bluebell/dao/redis"
	"bluebell.com/bluebell/models"
	"go.uber.org/zap"
	"strconv"
)

// 投票功能

// VoteForPost() 为帖子投票的函数

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
	// 1、判断投票限制
	// 2、更新分数
	// 3、记录用户为该帖子投票的数据
}
