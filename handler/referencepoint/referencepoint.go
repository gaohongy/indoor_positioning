package referencepoint

import (
	"indoor_positioning/model"
	"time"
)

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

type PutRequest struct {
	Id           uint64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Coordinate_x float64 `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64 `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64 `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
}

type PutResponse struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate   time.Time `json:"updatedate" gorm:"column:updatedate"`
}
