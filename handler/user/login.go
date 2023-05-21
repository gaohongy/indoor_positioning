package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/auth"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/zxmrlc/log"

	"github.com/gin-gonic/gin"
)

// @title	Login
// @description	用户登录API
// @auth	高宏宇
// @param	ctx *gin.Context
func Login(ctx *gin.Context) {
	log.Info("User Login function called")

	// 解析body参数
	var paraUser model.User
	if err := ctx.Bind(&paraUser); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// 查询用户
	dbUser, err := model.GetUserByUsername(paraUser.Username)
	if err != nil {
		log.Error("username error", err)
		// log.Infof("user (%s) is not found in database (%s)", paraUser.Username, viper.GetString("db.name"))
		handler.SendResponse(ctx, errno.ErrorLogin, nil)
		return
	}

	// 校验用户密码
	if err := auth.Compare(dbUser.Password, paraUser.Password); err != nil {
		log.Error("password error", err)
		handler.SendResponse(ctx, errno.ErrorLogin, nil)
		return
	}

	// 签发令牌
	token, err := token.Sign(token.Context{ID: dbUser.Id}, "")
	if err != nil {
		handler.SendResponse(ctx, errno.ErrorToken, nil)
		return
	}

	loginResponse := LoginResponse{
		UserType: dbUser.Usertype,
		Place_id: dbUser.Place_id,
		Token:    token,
	}

	handler.SendResponse(ctx, nil, loginResponse)
}
