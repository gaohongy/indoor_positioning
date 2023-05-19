package ap

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
// @description	查询符合时间条件的接入点API
// @auth	高宏宇
// @param	ctx *gin.Context
func Get(ctx *gin.Context) {
	log.Info("AP Get function called")

	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 解析请求数据
	begin_time_encode := ctx.Query("begin_time") // begin_time_encode: url安全格式编码
	end_time_encode := ctx.Query("end_time")     // // end_time_encode: url安全格式编码

	var ap_list_origin *[]model.Ap

	// 有时间筛选条件
	if begin_time_encode != "" && end_time_encode != "" {

		begin_time_decode, _ := url.QueryUnescape(begin_time_encode) // begin_time_decode: 原始时间字符串
		end_time_decode, _ := url.QueryUnescape(end_time_encode)     // end_time_decode: 原始时间字符串
		// 将字符串的时间转换为time.Time类型
		begin_time, _ := time.Parse(time.RFC3339, begin_time_decode)
		end_time, _ := time.Parse(time.RFC3339, end_time_decode)

		var err error
		// 查询接入点
		ap_list_origin, err = model.FilterApByTime(int(place_id), begin_time, end_time)
		if err != nil {
			log.Error("search ap_list_origin error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
		// 无时间筛选条件
	} else {
		// 请求ap_list数据
		var err error
		ap_list_origin, err = model.GetApByPlaceId(place_id)
		if err != nil {
			log.Error("search ap_list_origin error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
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
