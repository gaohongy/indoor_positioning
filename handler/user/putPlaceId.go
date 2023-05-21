package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	PutPlaceId
// @description	登录用户修改自身场所API
// @auth	高宏宇
// @param	ctx *gin.Context
func PutPlaceId(ctx *gin.Context) {
	log.Info("User Put PlaceId function called")

	// 解析body参数
	var request PutPlaceIdRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)

	// 修改用户所在场所
	if err := user.UpdateUserPlaceId(request.Place_id); err != nil {
		log.Error("user update place_id error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 发送响应
	handler.SendResponse(ctx, nil, nil)
}
