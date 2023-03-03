package gridpoint

import "indoor_positioning/model"

// place相关api所需请求响应结构
type CreateRequest struct {
	model.Coordinate
}

type CreateResponse struct {
	Gridpoint_id uint64 `json:"gridpoint_id"`
}
