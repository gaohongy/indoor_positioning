package location

import "indoor_positioning/model"

type Fingerprint struct {
	Bssid string  `json:"bssid"`
	Rss   float64 `json:"rss"`
}

// location相关api所需请求响应结构
type GetRequest struct {
	model.Coordinate
}

type GetResponse struct {
	model.Coordinate
}

type CreateKnnRequest struct {
	Offline_rss      [][]float64 `json:"offline_rss"`
	Offline_location [][]float64 `json:"offline_location"`
	Online_rss       [][]float64 `json:"online_rss"`
}
type CreateKnnResponse struct {
	Coordinate [][]float64 `json:"coordinate"`
}
