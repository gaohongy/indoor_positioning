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

// Create creates a new user account.
func Put(ctx *gin.Context) {
	log.Info("Ap Put function called")

	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)
	place_id := user.Place_id

	// 解析请求body
	var request PutRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(errno.ErrorBind.Error(), err)
		handler.SendResponse(ctx, errno.ErrorBind, nil)
		return
	}
	ap, _ := model.GetApById(request.Id)

	// 处理参考点
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

		// TODO 网格点如果插入失败，这里直接return是否可以
		if err := gridpoint.Create(); err != nil {
			log.Error("gridpoint insert error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		// 新插入网格点需要获取id
		gridpoint.Id = gridpoint.GetId()
	}

	if err := ap.Update(request.Ssid, request.Bssid, gridpoint.Id); err != nil {
		log.Error("user update place_id error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

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

	// TODO 验证参数合法性
	// if err := user.Validate(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorValidation, nil)
	// 	return
	// }

	// createResponse := CreateResponse{
	// 	Username: request.Username,
	// }

	// 发送响应
	handler.SendResponse(ctx, nil, putResopnse)
}
