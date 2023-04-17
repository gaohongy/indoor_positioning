package model

import (
	"time"
)

type Place struct {
	Id            uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Place_address string    `json:"place_address" gorm:"column:place_address;not null" binding:"required"`
	Longitude     float64   `json:"longitude" gorm:"column:longitude;not null" binding:"required"`
	Latitude      float64   `json:"latitude" gorm:"column:latitude;not null" binding:"required"`
	Createdate    time.Time `gorm:"column:createdate"`
	Updatedate    time.Time `gorm:"column:updatedate"`
}

type Place_brief struct {
	Place_address string  `json:"place_address" gorm:"column:place_address;not null" binding:"required"`
	Longitude     float64 `json:"longitude" gorm:"column:longitude;not null" binding:"required"`
	Latitude      float64 `json:"latitude" gorm:"column:latitude;not null" binding:"required"`
}

// 向数据库插入场所
func (place *Place) Create() error {
	return DB.Mysql.Create(&place).Error
}

func (place *Place) GetId() uint64 {
	t := &Place{}
	// TODO 添加查询失败时的处理
	DB.Mysql.Where("place_address = ?", place.Place_address).Find(&t)
	return t.Id

}

func GetAllPlaces() ([]Place, error) {
	places := []Place{}

	db := DB.Mysql.Find(&places)

	return places, db.Error
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
