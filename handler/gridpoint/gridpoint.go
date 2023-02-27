package gridpoint

// place相关api所需请求响应结构
type CreateRequest struct {
	Coordinate_x float64 `json:"coordinate_x"`
	Coordinate_y float64 `json:"coordinate_y"`
	Coordinate_z int     `json:"coordinate_z"`
}

type CreateResponse struct {
	Id uint64 `json:"id"`
}
