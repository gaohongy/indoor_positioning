package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	GetInfo
// @description	查询用户信息API
// @auth	高宏宇
// @param	ctx *gin.Context
func GetInfo(ctx *gin.Context) {
	log.Info("User Get Info function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)

	// 构造简短用户信息
	user_brief := model.User_Brief{
		Id:         user.Id,
		Username:   user.Username,
		Usertype:   user.Usertype,
		Place_id:   user.Place_id,
		Createdate: user.Createdate,
		Updatedate: user.Updatedate,
	}

	handler.SendResponse(ctx, nil, user_brief)
}
