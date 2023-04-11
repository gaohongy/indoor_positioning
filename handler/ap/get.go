package ap

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

func Get(ctx *gin.Context) {
	log.Info("AP Get function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

	place_id := user.Place_id

	// 请求ap_list数据
	ap_list_origin, err := model.GetApByPlaceId(int(place_id))
	if err != nil {
		// TODO 写入日志错误内容细化
		log.Error("search ap_list_origin error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 构造带有具体坐标的ap_list
	var ap_detail_list []model.Ap_Detail
	for _, ap_origin := range *ap_list_origin {

		// 查询所在网格点，以便获取坐标
		gridpoint, err := model.GetGridpointById(ap_origin.Grid_point_id)
		if err != nil {
			log.Error("search gridpoint error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		ap_detail := model.Ap_Detail{
			Id:           ap_origin.Id,
			Ssid:         ap_origin.Ssid,
			Bssid:        ap_origin.Bssid,
			Coordinate_x: gridpoint.Coordinate_x,
			Coordinate_y: gridpoint.Coordinate_y,
			Coordinate_z: gridpoint.Coordinate_z,
			Createdate:   ap_origin.Createdate,
			Updatedate:   ap_origin.Updatedate,
		}
		ap_detail_list = append(ap_detail_list, ap_detail)
	}
	handler.SendResponse(ctx, nil, ap_detail_list)
}
