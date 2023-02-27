package referencepoint

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
	log.Info("Referencepoint Create function called")

	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// TODO 需要根据用户token解析出place_id
	place_id := uint64(7)

	// 创建参考点前，其所在的网格点不一定存在。查询所在网格点，当网格点不存在时，新建网格点
	gridpoint, err := model.GetGridpoint(request.Coordinate_x, request.Coordinate_y, request.Coordinate_z, place_id)
	// 查询结果为空err.Error() = "record not found"
	if err != nil {
		log.Error("gridpoint not exists, creating a new gridpoint", err)

		gridpoint = &model.Gridpoint{
			Coordinate_x: request.Coordinate_x,
			Coordinate_y: request.Coordinate_y,
			Coordinate_z: request.Coordinate_z,
			Place_id:     place_id,
			Createdate:   time.Now(),
			Updatedate:   time.Now(),
		}

		// TODO 网格点如果插入失败，这里直接return是否可以
		if err := gridpoint.Create(); err != nil {
			log.Error(errno.ErrorDatabase.Error(), err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		// 新插入网格点需要获取id
		gridpoint.Id = gridpoint.GetId()
	}

	referencepoint := model.Referencepoint{
		Grid_point_id: gridpoint.Id,
		Place_id:      place_id,
		Createdate:    time.Now(),
		Updatedate:    time.Now(),
	}

	// TODO 验证参数合法性
	// if err := place.Validate(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorValidation, nil)
	// 	return
	// }

	// 参考点数据插入数据库
	if err := referencepoint.Create(); err != nil {
		log.Error(errno.ErrorDatabase.Error(), err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	createResponse := CreateResponse{
		Id: referencepoint.GetId(),
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
