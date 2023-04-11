package model

import (
	"time"
)

type Gridpoint struct {
	Id           uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Place_id     uint64    `json:"-"`
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

// 根据场所id，x，y，z坐标查询某场所的参考点
func GetGridpoint(coordinate_x float64, coordinate_y float64, coordinate_z float64, place_id uint64) (*Gridpoint, error) {
	gridpoint := &Gridpoint{}
	// 查询结果为空报错record not found
	db := DB.Mysql.Where("coordinate_x = ? AND coordinate_y = ? AND coordinate_z = ? AND place_id = ?",
		coordinate_x, coordinate_y, coordinate_z, place_id).Find(&gridpoint)
	return gridpoint, db.Error
}

func GetGridpointById(id uint64) (*Gridpoint, error) {
	// QUES 为何这里采用指针，可能是因为传参用指针更节省资源
	t := &Gridpoint{}
	// TODO 添加查询失败时的处理
	db := DB.Mysql.Where("id = ?", id).Find(&t)
	return t, db.Error
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
