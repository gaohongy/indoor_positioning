package gridpoint

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
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

	gridpoint := model.Gridpoint{
		Coordinate_x: request.Coordinate_x,
		Coordinate_y: request.Coordinate_y,
		Coordinate_z: request.Coordinate_z,
		// TODO 需要根据用户token解析出place_id
		Place_id:   7,
		Createdate: time.Now(),
		Updatedate: time.Now(),
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
		Id: gridpoint.GetId(),
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
