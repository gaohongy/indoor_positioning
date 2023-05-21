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

// @title	Create
// @description	新建网格点
// @auth	高宏宇
// @param	ctx *gin.Context
func Create(ctx *gin.Context) {
	log.Info("Gridpoint Create function called")

	// 解析body参数
	var request CreateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	gridpoint := model.Gridpoint{
		Coordinate_x: request.Coordinate_x,
		Coordinate_y: request.Coordinate_y,
		Coordinate_z: request.Coordinate_z,
		Place_id:     place_id,
		Createdate:   time.Now(),
		Updatedate:   time.Now(),
	}

	// 创建网格点
	if err := gridpoint.Create(); err != nil {
		log.Error(errno.ErrorDatabase.Error(), err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 响应数据
	createResponse := CreateResponse{
		Gridpoint_id: gridpoint.GetId(),
	}

	// 发送响应
	handler.SendResponse(ctx, nil, createResponse)
}
