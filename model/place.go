package model

import (
	"time"
)

// 场所信息
type Place struct {
	Id            uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`                         // 场所ID
	Place_address string    `json:"place_address" gorm:"column:place_address;not null" binding:"required"` // 场所详细地址
	Longitude     float64   `json:"longitude" gorm:"column:longitude;not null" binding:"required"`         // 场所所在经度
	Latitude      float64   `json:"latitude" gorm:"column:latitude;not null" binding:"required"`           // 场所所在纬度
	Createdate    time.Time `gorm:"column:createdate"`                                                     // 创建时间
	Updatedate    time.Time `gorm:"column:updatedate"`                                                     // 修改时间
}

// 场所信息（简要信息）
type Place_brief struct {
	Id            uint64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`                        // 场所ID
	Place_address string  `json:"place_address" gorm:"column:place_address;not null" binding:"required"` // 场所详细地址
	Longitude     float64 `json:"longitude" gorm:"column:longitude;not null" binding:"required"`         // 场所所在经度
	Latitude      float64 `json:"latitude" gorm:"column:latitude;not null" binding:"required"`           // 场所所在纬度
}

// @title	Create
// @description	结构体方法，向数据库插入场所
// @auth	高宏宇
// @receiver	place *Place 场所结构体对象指针
// @return	error 错误信息
func (place *Place) Create() error {
	return DB.Mysql.Create(&place).Error
}

// @title	GetId
// @description	结构体方法，查询场所ID
// @auth	高宏宇
// @receiver	place *Place 场所结构体对象指针
// @return	uint64	场所ID
func (place *Place) GetId() uint64 {
	t := &Place{}
	DB.Mysql.Where("place_address = ?", place.Place_address).Find(&t)
	return t.Id

}

// @title	GetAllPlaces
// @description	查询全部场所
// @auth	高宏宇
// @return	[]Place场所对象数组    error 错误信息
func GetAllPlaces() ([]Place, error) {
	places := []Place{}

	db := DB.Mysql.Find(&places)

	return places, db.Error
}
