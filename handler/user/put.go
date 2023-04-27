package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// Create creates a new user account.
func Put(ctx *gin.Context) {
	log.Info("User Put PlaceId function called")

	var request PutRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	user, _ := model.GetUserById(request.Id)

	if err := user.UpdateUserInfo(request.Username, request.Usertype); err != nil {
		log.Error("user update information error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	user, _ = model.GetUserById(request.Id)
	user_brief := model.User_Brief{
		Id:         user.Id,
		Username:   user.Username,
		Usertype:   user.Usertype,
		Place_id:   user.Place_id,
		Createdate: user.Createdate,
		Updatedate: user.Updatedate,
	}

	// TODO 验证参数合法性
	// if err := user.Validate(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorValidation, nil)
	// 	return
	// }

	// createResponse := CreateResponse{
	// 	Username: request.Username,
	// }

	// 发送响应
	handler.SendResponse(ctx, nil, user_brief)
}
