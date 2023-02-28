package gridpoint

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// Create creates a new user account.
func Create(ctx *gin.Context) {
	log.Info("Gridpoint Create function called")

	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

	place_id := user.Place_id

	gridpoint := model.Gridpoint{
		Coordinate_x: request.Coordinate_x,
		Coordinate_y: request.Coordinate_y,
		Coordinate_z: request.Coordinate_z,
		Place_id:     place_id,
		Createdate:   time.Now(),
		Updatedate:   time.Now(),
	}

	// TODO 验证参数合法性
	// if err := place.Validate(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorValidation, nil)
	// 	return
	// }

	// 场所数据插入数据库
	if err := gridpoint.Create(); err != nil {
		log.Error(errno.ErrorDatabase.Error(), err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	createResponse := CreateResponse{
		Gridpoint_id: gridpoint.GetId(),
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
