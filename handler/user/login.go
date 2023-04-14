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

// Create creates a new user account.
func Login(ctx *gin.Context) {
	log.Info("User Login function called")

	// Binding the data with the user struct.
	var paraUser model.User
	if err := ctx.Bind(&paraUser); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// Get the user information by the login username.
	dbUser, err := model.GetUserByUsername(paraUser.Username)
	if err != nil {
		// 如果需要写日志，需要写清详细信息-用户不存在，反馈给前端信息不能太详细
		log.Error("username error", err)
		// log.Infof("user (%s) is not found in database (%s)", paraUser.Username, viper.GetString("db.name"))
		handler.SendResponse(ctx, errno.ErrorLogin, nil)
		return
	}

	// Compare the login password with the user password.
	if err := auth.Compare(dbUser.Password, paraUser.Password); err != nil {
		log.Error("password error", err)
		handler.SendResponse(ctx, errno.ErrorLogin, nil)
		return
	}

	// Sign the json web token.
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
	// TODO 返回数据中需要添加用户类型，以便app端可以选择跳转页面
	handler.SendResponse(ctx, nil, loginResponse)
}
