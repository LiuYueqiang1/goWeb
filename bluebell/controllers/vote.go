package controllers

import (
	"bluebell.com/bluebell/logic"
	"bluebell.com/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 投票
//type VoteData struct {
//	PostID    int64 `json:"post_id,string"`
//	Direction int   `json:"direction,string"`
//}

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	logic.PostVote()
	ResponseSuccess(c, nil)
}
