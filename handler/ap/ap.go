// 接入点相关API所需请求响应结构
package ap

import (
	"indoor_positioning/model"
	"time"
)

// 创建接入点API请求结构
type CreateRequest struct {
	Ssid  string `json:"ssid"`
	Bssid string `json:"bssid"`
	model.Coordinate
}

// 创建接入点API响应结构
type CreateResponse struct {
	Ap_id uint64 `json:"ap_id"`
}

// 修改接入点API请求结构
type PutRequest struct {
	Id           uint64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Ssid         string  `json:"ssid" gorm:"column:ssid;not null" binding:"required"`
	Bssid        string  `json:"bssid" gorm:"column:bssid;not null" binding:"required"`
	Coordinate_x float64 `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64 `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64 `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
}

// 修改接入点API响应结构
type PutResponse struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Ssid         string    `json:"ssid" gorm:"column:ssid;not null" binding:"required"`
	Bssid        string    `json:"bssid" gorm:"column:bssid;not null" binding:"required"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate   time.Time `json:"updatedate" gorm:"column:updatedate"`
}
