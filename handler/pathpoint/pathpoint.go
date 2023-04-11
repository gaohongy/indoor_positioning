package pathpoint

import (
	"indoor_positioning/model"
)

// place相关api所需请求响应结构
type CreateRequest struct {
	model.Coordinate
}

// 前端传过来的时间值就是字符串，需要手动转换为time.Time格式
type GetRequest struct {
	Begin_time         string `json:"begin_time"`
	End_time           string `json:"end_time"`
	Selected_user_list string `json:"selected_user_list"`
}

type GetResponse struct {
	Pathpoint_list map[int][]model.Pathpoint_Detail `json:"pathpoint_list"`
}

// type CreateResponse struct {
// 	Id uint64 `json:"id"`
// }
