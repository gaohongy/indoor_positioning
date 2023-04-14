// TODO 目前按照每次刷新位置都会上传一次，但从实际角度考虑这样太耗时，可以考虑本地缓存一定量后一并上传
package pathpoint

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
	log.Info("Pathpoint Create function called")

	// 解析请求参数
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

	// 创建路径点前，其所在的网格点不一定存在。查询所在网格点，当网格点不存在时，新建网格点
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
			log.Error("gridpoint insert error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		// 新插入网格点需要获取id
		// TODO 可以直接利用插入点时的返回值，可以少查询一次
		gridpoint.Id = gridpoint.GetId()
	}

	pathpoint := model.Pathpoint{
		User_id:       user.Id,
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

	// 路径点数据插入数据库
	if err := pathpoint.Create(); err != nil {
		log.Error("pathpoint insert error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// createResponse := CreateResponse{
	// 	Id: gridpoint.GetId(),
	// }

	// 发送响应
	handler.SendResponse(ctx, nil, nil)
}
