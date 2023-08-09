package redis

// redis key

//redis key 注意使用命名空间的方式，方便拆分和查询

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"   // zset；帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  //zset;帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset;记录用户及投票类型;参数是 post id

)
