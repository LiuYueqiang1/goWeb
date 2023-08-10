package logic

import (
	"bluebell.com/bluebell/dao/mysql"
	"bluebell.com/bluebell/dao/redis"
	"bluebell.com/bluebell/models"
	"bluebell.com/bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
	// 3. 返回
}

// GetPostById1 根据帖子id查询帖子基本信息
/*
   "id": "7290692015493120",
   "author_id": 0,
   "community_id": 1,
   "status": 0,
   "title": "study go make me happy",
   "content": "My ei must be accept",
   "create_time": "2023-07-21T10:50:36Z"
*/
func GetPostById1(pid int64) (data *models.Post, err error) {
	return mysql.GetPostById(pid)
}

// GetPostById2 根据帖子id查询帖子详情数据
// 根据上面做改造
/*
   "author_name": "七米",
   "id": "7367066386436096",
   "author_id": 1837296026390528,
   "community_id": 2,
   "status": 0,
   "title": "study make me happy",
   "content": "全心全意攻占浙大学报",
   "create_time": "2023-07-21T15:54:05Z",
   "community": {
       "id": 2,
       "name": "leetcode",
       "introduction": "刷题刷题",
       "create_time": "2023-07-01T08:10:10Z"
*/
func GetPostById2(pid int64) (data *models.ApiPostDetail, err error) {

	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	// 必须要从中间件 post 的信息才可以，这样才有AuthorID
	// 不然拿不到 AuthorID 的话将会返回错误
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.AuthorID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	////// 投票数
	////voteData, err := redis.GetPostVoteData([]string{strconv.Itoa(int(pid))})
	////if err != nil {
	////	return
	////}
	// 接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
		//VoteNum:         voteData[0],
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
