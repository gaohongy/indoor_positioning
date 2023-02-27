package ap

// 人工添加AP调用
// ap相关api所需请求响应结构
type CreateRequest struct {
	Ssid         string  `json:"ssid"`
	Bssid        string  `json:"bssid"`
	Coordinate_x float64 `json:"coordinate_x"`
	Coordinate_y float64 `json:"coordinate_y"`
	Coordinate_z int     `json:"coordinate_z"`
}

type CreateResponse struct {
	Id uint64 `json:"id"`
}
