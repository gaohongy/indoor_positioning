package referencepoint

import "indoor_positioning/model"

type fingerPrint struct {
	Bssid string  `json:"bssid"`
	Rss   float64 `json:"rss"`
}

// 这里的调用场景是：安卓端在添加参考点时，扫描到很多wifi信息，同时需要手动输入x、y、z，然后提交。所以表单中能接收到的数据就只有x、y、z，所需的场所id应当从token中进行解析
// 在搜集rss信息时，会输入当前位置x，y，z，创建时会首先创建参考点，然后返回参考点id，用于创建rss
type CreateRequest struct {
	Rss_list []fingerPrint `json:"rss_list"`
	model.Coordinate
}

type CreateResponse struct {
	Referencepoint_id uint64 `json:"referencepoint_id"`
}
