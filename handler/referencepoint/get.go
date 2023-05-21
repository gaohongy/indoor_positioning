package referencepoint

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Get
// @description	查询符合时间条件的参考点API
// @auth	高宏宇
// @param	ctx *gin.Context
func Get(ctx *gin.Context) {
	log.Info("referencepoint Get function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 解析请求数据
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")

	var referencepoint_list_origin []*model.Referencepoint

	// 有时间筛选条件
	if begin_time_encode != "" && end_time_encode != "" {
		// 将url安全格式编码字符串还原为原始字符串
		begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
		end_time_decode, _ := url.QueryUnescape(end_time_encode)
		// 将字符串的时间转换为time.Time类型
		begin_time, _ := time.Parse(time.RFC3339, begin_time_decode)
		end_time, _ := time.Parse(time.RFC3339, end_time_decode)

		// 带筛选条件查询参考点
		var err error
		referencepoint_list_origin, err = model.FilterReferencepointByTime(int(place_id), begin_time, end_time)
		if err != nil {
			log.Error("search ap_list_origin error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	} else {
		// 查询当前场所所有的参考点
		var err error
		referencepoint_list_origin, _, err = model.ListReferencepointByPlaceid(place_id, 0, 0)
		if err != nil {
			log.Error("search referencepoint list error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	}

	// 构造带有具体坐标的参考点列表
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
