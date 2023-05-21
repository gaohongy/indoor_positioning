package location

import "indoor_positioning/model"

type Fingerprint struct {
	Bssid string  `json:"bssid"`
	Rss   float64 `json:"rss"`
}

// 预测坐标请求数据结构
type GetRequest struct {
	model.Coordinate
}

// 预测坐标响应数据结构
type GetResponse struct {
	model.Coordinate
}

// WKNN坐标计算请求数据结构
type CreateKnnRequest struct {
	Offline_rss      [][]float64 `json:"offline_rss"`      // 离线数据-信号强度
	Offline_location [][]float64 `json:"offline_location"` // 离线数据-参考点坐标
	Online_rss       [][]float64 `json:"online_rss"`       // 在线数据-信号强度
}

// WKNN坐标计算响应数据结构
type CreateKnnResponse struct {
	Coordinate [][]float64 `json:"coordinate"`
}
