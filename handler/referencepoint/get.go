package referencepoint

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

func Get(ctx *gin.Context) {
	log.Info("referencepoint Get function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)

	place_id := user.Place_id

	// 查询当前场所所有的参考点
	referencepoint_list_origin, _, err := model.ListReferencepointByPlaceid(place_id, 0, 0)
	if err != nil {
		log.Error("search referencepoint list error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 构造带有具体坐标的referencepoint_list
	var referencepoint_detail_list []model.Referencepoint_Detail
	for _, referencepoint_origin := range referencepoint_list_origin {

		// 查询所在网格点，以便获取坐标
		gridpoint, err := model.GetGridpointById(referencepoint_origin.Grid_point_id)
		if err != nil {
			log.Error("search gridpoint error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}

		referencepoint_detail := model.Referencepoint_Detail{
			Id:           referencepoint_origin.Id,
			Coordinate_x: gridpoint.Coordinate_x,
			Coordinate_y: gridpoint.Coordinate_y,
			Coordinate_z: gridpoint.Coordinate_z,
			Createdate:   referencepoint_origin.Createdate,
			Updatedate:   referencepoint_origin.Updatedate,
		}
		referencepoint_detail_list = append(referencepoint_detail_list, referencepoint_detail)
	}
	handler.SendResponse(ctx, nil, referencepoint_detail_list)
}
