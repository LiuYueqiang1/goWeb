package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//定义请求的参数结构体

// ParmSignUp 注册用户的结构体
type ParmSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	// "required,eqfield=Password" 判断re_password == password
}

// ParmLogin 用户登录的结构体
type ParmLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required"`
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票+1 反对票-1  取消投票0
}

// ParamPostList 获取帖子列表 string 的参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空
	Page        int64  `json:"page" form:"page"`                 // 页码
	Size        int64  `json:"size" form:"size"`                 // 每页数量
	Order       string `json:"order" form:"order"`               // 排序依据
}

// ParamCommunityPostList 按社区获取帖子列表query string 参数

type ParamCommunityPostList struct {
	*ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
