package gridpoint

import "indoor_positioning/model"

// 创建网格点请求结构
type CreateRequest struct {
	model.Coordinate
}

// 创建网格点响应结构
type CreateResponse struct {
	Gridpoint_id uint64 `json:"gridpoint_id"`
}
