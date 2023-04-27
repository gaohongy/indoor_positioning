package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// Create creates a new user account.
func PutPlaceId(ctx *gin.Context) {
	log.Info("User Put PlaceId function called")

	var request PutRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

	if err := user.Update(request.Place_id); err != nil {
		log.Error("user update place_id error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
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
	handler.SendResponse(ctx, nil, nil)
}
