package place

// place相关api所需请求响应结构
type CreateRequest struct {
	Place_address string  `json:"place_address"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
}

type CreateResponse struct {
	Id uint64 `json:"id"`
}
