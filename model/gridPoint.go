package model

import (
	"time"
)

type Gridpoint struct {
	Id           uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z int       `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Place_id     int       `json:"-"`
	Createdate   time.Time `gorm:"column:createdate"`
	Updatedate   time.Time `gorm:"column:updatedate"`
}

// 向数据库插入场所
func (gridpoint *Gridpoint) Create() error {
	return DB.Mysql.Create(&gridpoint).Error
}

func (gridpoint *Gridpoint) GetId() uint64 {
	t := &Gridpoint{}
	// TODO 添加查询失败时的处理
	DB.Mysql.Where("coordinate_x = ? AND coordinate_y = ? AND coordinate_z = ?",
		gridpoint.Coordinate_x, gridpoint.Coordinate_y, gridpoint.Coordinate_z).Find(&t)
	return t.Id
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
