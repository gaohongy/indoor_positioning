package ap

import "indoor_positioning/model"

// 人工添加AP调用
// ap相关api所需请求响应结构
type CreateRequest struct {
	Ssid  string `json:"ssid"`
	Bssid string `json:"bssid"`
	model.Coordinate
}

type CreateResponse struct {
	Ap_id uint64 `json:"ap_id"`
}
