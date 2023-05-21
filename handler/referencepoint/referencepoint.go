package referencepoint

import (
	"indoor_positioning/model"
	"time"
)

// WiFi指纹信息
type fingerPrint struct {
	Bssid string  `json:"bssid"` // 基本服务集标识符
	Rss   float64 `json:"rss"`   // 接收信号强度
}

// 创建参考点API请求数据结构
type CreateRequest struct {
	Rss_list []fingerPrint `json:"rss_list"`
	model.Coordinate
}

// 创建参考点API响应数据结构
type CreateResponse struct {
	Referencepoint_id uint64 `json:"referencepoint_id"`
}

// 修改参考点API请求数据结构
type PutRequest struct {
	Id           uint64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Coordinate_x float64 `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64 `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64 `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
}

// 修改参考点API响应数据结构
type PutResponse struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate   time.Time `json:"updatedate" gorm:"column:updatedate"`
}
