package pathpoint

import (
	"indoor_positioning/model"
)

// 创建路径点请求数据结构
type CreateRequest struct {
	model.Coordinate
}

// 筛选指定时间段内指定用户历史路径请求数据结构
type GetRequest struct {
	Begin_time         string `json:"begin_time"`
	End_time           string `json:"end_time"`
	Selected_user_list string `json:"selected_user_list"`
}

// 筛选指定时间段内指定用户历史路径响应数据结构
type GetResponse struct {
	Pathpoint_list map[int][]model.Pathpoint_Detail `json:"pathpoint_list"`
}
