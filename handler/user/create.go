package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// Create creates a new user account.
func Create(ctx *gin.Context) {
	log.Info("User Create function called.")

	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// TODO 能否在创建对象时自动设置salt
	user := model.UserModel{
		Username: request.Username,
		Password: request.Password,
	}

	// 验证参数合法性
	if err := user.Validate(); err != nil {
		handler.SendResponse(ctx, errno.ErrorValidation, nil)
		return
	}

	// 加密用户密码
	// if err := user.Encrypt(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorEncrypt, nil)
	// 	return
	// }

	// 用户数据插入数据库
	if err := user.Create(); err != nil {
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	createResponse := CreateResponse{
		Username: request.Username,
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
