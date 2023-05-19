package ap

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Get
// @description	修改接入点信息API
// @auth	高宏宇
// @param	ctx *gin.Context
func Put(ctx *gin.Context) {
	log.Info("Ap Put function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 解析body参数
	var request PutRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}
	// 查询待修改接入点
	ap, _ := model.GetApById(request.Id)

	// 查询参考点
	gridpoint, err := model.GetGridpoint(request.Coordinate_x, request.Coordinate_y, request.Coordinate_z, place_id)
	// 查询结果为空err.Error() = "record not found"
	if err != nil {
		log.Info("gridpoint not exists, creating a new gridpoint")

		gridpoint = &model.Gridpoint{
			Coordinate_x: request.Coordinate_x,
			Coordinate_y: request.Coordinate_y,
			Coordinate_z: request.Coordinate_z,
			Place_id:     place_id,
			Createdate:   time.Now(),
			Updatedate:   time.Now(),
		}

		if err := gridpoint.Create(); err != nil {
			log.Error("gridpoint insert error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		// 新插入网格点需要获取id
		gridpoint.Id = gridpoint.GetId()
	}

	// 更新接入点信息
	if err := ap.Update(request.Ssid, request.Bssid, gridpoint.Id); err != nil {
		log.Error("AP update place_id error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 查询更新后的接入点
	ap, _ = model.GetApById(request.Id)

	putResopnse := PutResponse{
		Id:           ap.Id,
		Ssid:         ap.Ssid,
		Bssid:        ap.Bssid,
		Coordinate_x: gridpoint.Coordinate_x,
		Coordinate_y: gridpoint.Coordinate_y,
		Coordinate_z: gridpoint.Coordinate_z,
		Createdate:   ap.Createdate,
		Updatedate:   ap.Updatedate,
	}

	// 发送响应
	handler.SendResponse(ctx, nil, putResopnse)
}
