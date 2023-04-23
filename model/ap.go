package model

import (
	"time"

	"github.com/zxmrlc/log"
)

// `json:"-"`在查询结构映射为结构体时是不会映射的，也就是没有结果，但是createdate和updatedate似乎不太一样
type Ap struct {
	Id            uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Ssid          string    `json:"ssid" gorm:"column:ssid;not null" binding:"required"`
	Bssid         string    `json:"bssid" gorm:"column:bssid;not null" binding:"required"`
	Grid_point_id uint64    `json:"grid_point_id" gorm:"column:grid_point_id;not null"`
	Place_id      uint64    `json:"place_id" gorm:"column:place_id;not null"`
	Createdate    time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate    time.Time `json:"updatedate" gorm:"column:updatedate"`
}

type Ap_Detail struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Ssid         string    `json:"ssid" gorm:"column:ssid;not null" binding:"required"`
	Bssid        string    `json:"bssid" gorm:"column:bssid;not null" binding:"required"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate   time.Time `json:"updatedate" gorm:"column:updatedate"`
}

// 向数据库插入场所
func (ap *Ap) Create() error {
	return DB.Mysql.Create(&ap).Error
}

func (ap *Ap) GetId() uint64 {
	t := &Ap{}
	// TODO 添加查询失败时的处理
	db := DB.Mysql.Where("bssid = ?", ap.Bssid).Find(&t)
	if db.Error != nil {
		log.Error("ap.GetId() error", db.Error)
	}
	return t.Id
}

// 修改AP信息
func (ap *Ap) Update(ssid string, bssid string, grid_point_id uint64) error {
	db := DB.Mysql.Model(ap).Update(map[string]interface{}{"ssid": ssid, "bssid": bssid, "grid_point_id": grid_point_id})
	return db.Error
}

func GetApById(id uint64) (*Ap, error) {
	ap := &Ap{}
	db := DB.Mysql.Where("id = ?", id).Find(&ap)
	return ap, db.Error
}

func GetApByBssid(bssid string) (*Ap, error) {
	ap := &Ap{}
	db := DB.Mysql.Where("bssid = ?", bssid).Find(&ap)
	return ap, db.Error
}

func GetApByPlaceId(place_id int) (*[]Ap, error) {
	ap_list := &[]Ap{}
	db := DB.Mysql.Where("place_id = ?", place_id).Find(&ap_list)
	return ap_list, db.Error
}

func FilterApByTime(place_id int, begin_time time.Time, end_time time.Time) (*[]Ap, error) {
	ap_list := &[]Ap{}
	db := DB.Mysql.Where("place_id = ? AND createdate BETWEEN ? AND ?", place_id, begin_time, end_time).Find(&ap_list)
	return ap_list, db.Error
}

func DeleteAp(id uint64) error {
	ap := &Ap{}
	db := DB.Mysql.Where("id = ?", id).Delete(&ap)
	return db.Error
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
