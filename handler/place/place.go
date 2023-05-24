package place

// 创建场所请求数据结构
type CreateRequest struct {
	Place_address string  `json:"place_address"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	Map_id        string  `json:"map_id"`
}

// 创建场所响应数据结构
type CreateResponse struct {
	Place_id uint64 `json:"place_id"`
}
