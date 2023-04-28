package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

func GetInfo(ctx *gin.Context) {
	log.Info("User Get Info function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

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
