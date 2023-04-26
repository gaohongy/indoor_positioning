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

// 1. 参数解析
// 1.1 将url中的请求参数绑定到数据结构之中
// 1.2 将序列化的用户数组反序列化为数组
// 1.3 解析token获取场所id
// 2. 数据查询
// 3. 返回数据构造
// 4. 数据返回
// NOTE 有2种特殊情况：1.无时间筛选条件 2.无用户筛选条件 对于前者我们只需要在查询数据库时不添加时间这一判断即可，但是后者我们不能直接筛选时间，因为返回的数据需要根据用户进行分类以便于地图显示，所以必须一个用户一个用户来查询
func Get(ctx *gin.Context) {
	// TODO 改变user_id获取方式，或通过中间件实现
	content, _ := token.ParseRequest(ctx)
	user, _ := model.GetUserById(content.ID)
	place_id := user.Place_id

	// 解析请求数据
	begin_time_encode := ctx.Query("begin_time")
	end_time_encode := ctx.Query("end_time")
	selected_user_list_encode := ctx.Query("selected_user_list")

	// begin_time 和 end_time 此时为空，值为"0001-01-01 00:00:00 +0000 UTC"，begin_time.IsZero() = true
	var begin_time time.Time
	var end_time time.Time
	selected_user_list := []int{}

	// 反序列化参数中的用户数组
	if selected_user_list_encode != "" {
		selected_user_list_decode, err := url.QueryUnescape(selected_user_list_encode)
		if err != nil {
			// TODO 处理错误
			log.Error("QueryUnescape error", err)
		}
		json.Unmarshal([]byte(selected_user_list_decode), &selected_user_list)
		if err != nil {
			log.Error("selected_user_list unmarshal error", err)
			return
		}
	} else { // 没给出用户条件 <=> 所有用户
		all_user_list, err := model.GetUserByPlaceId(int(place_id))
		if err != nil {
			// TODO 处理错误
			log.Error("pathpoint get.go GetUserByPlaceId() error", err)
			return
		}

		for _, user := range all_user_list {
			selected_user_list = append(selected_user_list, int(user.Id))
		}
	}

	// 请求满足条件的路径点数据
	var response GetResponse
	// NOTE map赋值前要先初始化，否则会报错 assignment to entry in nil map
	response.Pathpoint_list = make(map[int][]model.Pathpoint_Detail)
	// 一次查询一个用户
	for _, user_id := range selected_user_list {
		// 将字符串的时间转换为time.Time类型
		if begin_time_encode != "" && end_time_encode != "" {
			begin_time_decode, _ := url.QueryUnescape(begin_time_encode)
			end_time_decode, _ := url.QueryUnescape(end_time_encode)

			// TODO begin_time 和 end_time解析出来后都是UTC时间，但mysql中的时区未必是UTC，这里需要解决一下
			begin_time, _ = time.Parse(time.RFC3339, begin_time_decode)
			end_time, _ = time.Parse(time.RFC3339, end_time_decode)
		}

		// 根据场所id，用户id，开始时间和结束时间，搜索其中一个用户对应的所有路径点
		// 如果未添加时间筛选条件，那么begin_time和end_time此时为time.Time的空值
		pathpoint_origin_list, err := model.FilterPathpointByTimeAndUser(int(place_id), user_id, begin_time, end_time)
		if err != nil {
			// TODO 写入日志错误内容细化
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

	handler.SendResponse(ctx, nil, response)
}
