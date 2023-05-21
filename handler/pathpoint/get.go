package pathpoint

import (
	"encoding/json"
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
// @description	查询指定时间段内指定用户的历史轨迹
// @auth	高宏宇
// @param	ctx *gin.Context
func Get(ctx *gin.Context) {
	// 获取登录用户ID
	content, _ := token.ParseRequest(ctx)
	// 查询用户
	user, _ := model.GetUserById(content.ID)
	// 查询用户所在场所ID
	place_id := user.Place_id

	// 解析请求参数
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")
	selected_user_list_encode := ctx.Query("selected_user_list")

	// begin_time 和 end_time 此时为空，值为"0001-01-01 00:00:00 +0000 UTC"，begin_time.IsZero() = true
	var begin_time time.Time
	var end_time time.Time
	selected_user_list := []int{}

	// 用户筛选条件非空
	if selected_user_list_encode != "" {
		// 将url安全格式编码字符串还原为原始字符串
		selected_user_list_decode, err := url.QueryUnescape(selected_user_list_encode)
		if err != nil {
			log.Error("QueryUnescape error", err)
		}
		// 反序列化获得Go语言用户数组
		json.Unmarshal([]byte(selected_user_list_decode), &selected_user_list)
		if err != nil {
			log.Error("selected_user_list unmarshal error", err)
			return
		}
		// 未给出用户筛选条件 <=> 查询所有用户
	} else {
		// 查询当前场所全部用户
		all_user_list, err := model.GetUserByPlaceId(int(place_id))
		if err != nil {
			log.Error("pathpoint get.go GetUserByPlaceId() error", err)
			return
		}

		for _, user := range all_user_list {
			selected_user_list = append(selected_user_list, int(user.Id))
		}
	}

	// 请求满足条件的路径点数据
	var response GetResponse
	response.Pathpoint_list = make(map[int][]model.Pathpoint_Detail)
	// 一次查询一个用户
	for _, user_id := range selected_user_list {
		// 时间筛选条件非空
		if begin_time_encode != "" && end_time_encode != "" {
			// 将url安全格式编码字符串还原为原始字符串
			begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
			end_time_decode, _ := url.QueryUnescape(end_time_encode)

			// 将字符串的时间转换为time.Time类型
			begin_time, _ = time.Parse(time.RFC3339, begin_time_decode)
			end_time, _ = time.Parse(time.RFC3339, end_time_decode)
		}

		// 根据场所id，用户id，开始时间和结束时间，搜索其中一个用户对应的所有路径点
		// 如果未添加时间筛选条件，那么begin_time和end_time此时为time.Time的空值
		pathpoint_origin_list, err := model.FilterPathpointByTimeAndUser(int(place_id), user_id, begin_time, end_time)
		if err != nil {
			log.Error("filter pathpoint error", err)
			handler.SendResponse(ctx, errno.ErrorDatabase, nil)
			return
		}
		// 一个用户在筛选条件下不具有路径点则不返回数据
		if len(*pathpoint_origin_list) == 0 {
			continue
		}

		// 由于从数据库中直接拿到的数据并不包含x，y，z坐标，需要把上面拿到的数组中的路径点全部替换为包含具体坐标信息的路径点，从新构造一个数组
		var pathpoint_detail_list []model.Pathpoint_Detail
		for _, pathpoint_origin := range *pathpoint_origin_list {
			grid_point, _ := model.GetGridpointById(pathpoint_origin.Grid_point_id)

			pathpoint_detail := model.Pathpoint_Detail{
				Id:           pathpoint_origin.Id,
				Coordinate_x: grid_point.Coordinate_x,
				Coordinate_y: grid_point.Coordinate_y,
				Coordinate_z: grid_point.Coordinate_z,
				Createdate:   pathpoint_origin.Createdate,
			}

			pathpoint_detail_list = append(pathpoint_detail_list, pathpoint_detail)
		}

		// 把构造好的数组添加进map中
		response.Pathpoint_list[user_id] = pathpoint_detail_list
	}

	// 返回响应结果
	handler.SendResponse(ctx, nil, response)
}
