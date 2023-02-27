package referencepoint

// 在搜集rss信息时，会输入当前位置x，y，z，创建时会首先创建参考点，然后返回参考点id，用于创建rss
type CreateRequest struct {
	Coordinate_x float64 `json:"coordinate_x"`
	Coordinate_y float64 `json:"coordinate_y"`
	Coordinate_z int     `json:"coordinate_z"`
}

type CreateResponse struct {
	Id uint64 `json:"id"`
}
