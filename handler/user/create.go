package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Create
// @description	新建用户API
// @auth	高宏宇
// @param	ctx *gin.Context
func Create(ctx *gin.Context) {
	log.Info("User Create function called")

	// 解析body参数
	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// 检查是否存在同名用户
	_, err := model.GetUserByUsername(request.Username)
	// 用户已存在
	if err == nil {
		handler.SendResponse(ctx, errno.ErrorUsernameRepeat, nil)
		return
	}

	user := model.User{
		Username:   request.Username,
		Password:   request.Password,
		Usertype:   request.UserType,
		Createdate: time.Now(),
		Updatedate: time.Now(),
	}

	// 加密用户密码
	if err := user.Encrypt(); err != nil {
		handler.SendResponse(ctx, errno.ErrorEncrypt, nil)
		return
	}

	// 用户数据插入数据库
	if err := user.Create(); err != nil {
		log.Error(errno.ErrorDatabase.Error(), err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	createResponse := CreateResponse{
		Username: request.Username,
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
