package model

import (
	"time"

	"github.com/zxmrlc/log"
)

type Ap struct {
	Id            uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Ssid          string    `json:"ssid" gorm:"column:ssid;not null" binding:"required"`
	Bssid         string    `json:"bssid" gorm:"column:bssid;not null" binding:"required"`
	Grid_point_id uint64    `json:"-" gorm:"column:grid_point_id;not null"`
	Place_id      uint64    `json:"-" gorm:"column:place_id;not null"`
	Createdate    time.Time `gorm:"column:createdate"`
	Updatedate    time.Time `gorm:"column:updatedate"`
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

func GetAp(bssid string) (*Ap, error) {
	ap := &Ap{}
	db := DB.Mysql.Where("bssid = ?", bssid).Find(&ap)
	return ap, db.Error
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
