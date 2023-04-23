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

func Get(ctx *gin.Context) {
	log.Info("AP Get function called")

	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)
	place_id := user.Place_id

	// 解析请求数据
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")

	var ap_list_origin *[]model.Ap

	if begin_time_encode != "" && end_time_encode != "" { // 有时间筛选条件
		begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
		end_time_decode, _ := url.QueryUnescape(end_time_encode)
		// 将字符串的时间转换为time.Time类型
		// TODO begin_time 和 end_time解析出来后都是UTC时间，但mysql中的时区未必是UTC，这里需要解决一下
		begin_time, _ := time.Parse(time.RFC3339, begin_time_decode)
		end_time, _ := time.Parse(time.RFC3339, end_time_decode)

		var err error
		ap_list_origin, err = model.FilterApByTime(int(place_id), begin_time, end_time)
		if err != nil {
			// TODO 写入日志错误内容细化
			log.Error("search ap_list_origin error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
	} else { // 无时间筛选条件
		// 请求ap_list数据
		var err error
		ap_list_origin, err = model.GetApByPlaceId(int(place_id))
		if err != nil {
			// TODO 写入日志错误内容细化
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
