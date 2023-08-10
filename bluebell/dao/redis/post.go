package redis

import (
	"bluebell.com/bluebell/models"
	"context"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1、根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2、确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 3、ZREVRANGE 按分数从大到小查询指定数量的元素
	return rdb.ZRevRange(context.Background(), key, start, end).Result()
}
