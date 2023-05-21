package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Put
// @description	修改用户信息API
// @auth	高宏宇
// @param	ctx *gin.Context
func Put(ctx *gin.Context) {
	log.Info("User Put PlaceId function called")

	// 解析body参数
	var request PutRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// 查询用户信息
	user, _ := model.GetUserById(request.Id)

	// 更新用户信息
	if err := user.UpdateUserInfo(request.Username, request.Usertype); err != nil {
		log.Error("user update information error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 查询更新后用户信息
	user, _ = model.GetUserById(request.Id)
	user_brief := model.User_Brief{
		Id:         user.Id,
		Username:   user.Username,
		Usertype:   user.Usertype,
		Place_id:   user.Place_id,
		Createdate: user.Createdate,
		Updatedate: user.Updatedate,
	}

	// 发送响应
	handler.SendResponse(ctx, nil, user_brief)
}
