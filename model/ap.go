package model

import (
	"time"

	"github.com/zxmrlc/log"
)

// 接入点信息
type Ap struct {
	Id            uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`        // 接入点ID
	Ssid          string    `json:"ssid" gorm:"column:ssid;not null" binding:"required"`   // 服务集标识
	Bssid         string    `json:"bssid" gorm:"column:bssid;not null" binding:"required"` // 基本服务集标识
	Grid_point_id uint64    `json:"grid_point_id" gorm:"column:grid_point_id;not null"`    // 网格点ID
	Place_id      uint64    `json:"place_id" gorm:"column:place_id;not null"`              // 场所ID
	Createdate    time.Time `json:"createdate" gorm:"column:createdate"`                   // 创建时间
	Updatedate    time.Time `json:"updatedate" gorm:"column:updatedate"`                   // 修改时间
}

// 接入点信息（包含坐标详细信息）
type Ap_Detail struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`                      // 接入点ID
	Ssid         string    `json:"ssid" gorm:"column:ssid;not null" binding:"required"`                 // 服务集标识
	Bssid        string    `json:"bssid" gorm:"column:bssid;not null" binding:"required"`               // 基本服务集标识
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"` // 接入点所在x坐标
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"` // 接入点所在y坐标
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"` // 接入点所在z坐标
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`                                 // 创建时间
	Updatedate   time.Time `json:"updatedate" gorm:"column:updatedate"`                                 // 修改时间
}

// @title	Create
// @description	结构体方法，向数据库插入场所
// @auth	高宏宇
// @receiver	ap *Ap 接入点结构体对象指针
// @return	error	错误信息
func (ap *Ap) Create() error {
	return DB.Mysql.Create(&ap).Error
}

// @title	GetId
// @description	结构体方法，提供接入点BSSID查询接入点ID
// @auth	高宏宇
// @receiver	ap *Ap 接入点结构体对象指针
// @return	uint64	接入点ID
func (ap *Ap) GetId() uint64 {
	t := &Ap{}
	// TODO 添加查询失败时的处理
	db := DB.Mysql.Where("bssid = ?", ap.Bssid).Find(&t)
	if db.Error != nil {
		log.Error("ap.GetId() error", db.Error)
	}
	return t.Id
}

// @title	Update
// @description	结构体方法，修改AP信息
// @auth	高宏宇
// @receiver	ap *Ap 接入点结构体对象指针
// @param	ssid string 服务集标识	bssid string 基本服务集标识    grid_point_id uint64 网格点ID
// @return	error	错误信息
func (ap *Ap) Update(ssid string, bssid string, grid_point_id uint64) error {
	db := DB.Mysql.Model(ap).Update(map[string]interface{}{"ssid": ssid, "bssid": bssid, "grid_point_id": grid_point_id})
	return db.Error
}

// @title	GetApById
// @description	根据ID查询接入点
// @auth	高宏宇
// @param	id uint64 接入点ID
// @return	*Ap 接入点对象指针    error 错误信息
func GetApById(id uint64) (*Ap, error) {
	ap := &Ap{}
	db := DB.Mysql.Where("id = ?", id).Find(&ap)
	return ap, db.Error
}

// @title	GetApByBssid
// @description	根据BSSID查询接入点
// @auth	高宏宇
// @param	bssid string 基本服务集标识符
// @return	*Ap 接入点对象指针    error 错误信息
func GetApByBssid(bssid string) (*Ap, error) {
	ap := &Ap{}
	db := DB.Mysql.Where("bssid = ?", bssid).Find(&ap)
	return ap, db.Error
}

// @title	GetApByPlaceId
// @description	查询指定场所中的接入点
// @auth	高宏宇
// @param	place_id uint64 场所ID
// @return	*[]Ap 接入点对象数组指针    error 错误信息
func GetApByPlaceId(place_id uint64) (*[]Ap, error) {
	ap_list := &[]Ap{}
	db := DB.Mysql.Where("place_id = ?", place_id).Find(&ap_list)
	return ap_list, db.Error
}

// @title	FilterApByTime
// @description	筛选指定场所中指定时间段内添加的的接入点
// @auth	高宏宇
// @param	place_id uint64 场所ID    begin_time time.Time 筛选开始时间    end_time time.Time 筛选结束时间
// @return	*[]Ap 接入点对象数组指针    error 错误信息
func FilterApByTime(place_id int, begin_time time.Time, end_time time.Time) (*[]Ap, error) {
	ap_list := &[]Ap{}
	db := DB.Mysql.Where("place_id = ? AND createdate BETWEEN ? AND ?", place_id, begin_time, end_time).Find(&ap_list)
	return ap_list, db.Error
}

// @title	DeleteAp
// @description	删除接入点
// @auth	高宏宇
// @param	id uint64 接入点ID
// @return	error 错误信息
func DeleteAp(id uint64) error {
	ap := &Ap{}
	db := DB.Mysql.Where("id = ?", id).Delete(&ap)
	return db.Error
}
