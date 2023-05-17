package model

import (
	"time"
)

// 网格点信息
type Gridpoint struct {
	Id           uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`                       // 网格点ID
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"` // 网格点x坐标
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"` // 网格点y坐标
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"` // 网格点z坐标
	Place_id     uint64    `json:"-"`                                                                   // 场所ID
	Createdate   time.Time `gorm:"column:createdate"`                                                   // 创建时间
	Updatedate   time.Time `gorm:"column:updatedate"`                                                   // 修改时间
}

// @title	Create
// @description	结构体方法，向数据库插入场所
// @auth	高宏宇
// @receiver	gridpoint *Gridpoint 网格点结构体对象指针
// @return	error	错误信息
func (gridpoint *Gridpoint) Create() error {
	return DB.Mysql.Create(&gridpoint).Error
}

// @title	GetId
// @description	结构体方法，提供网格点坐标查询网格点ID
// @auth	高宏宇
// @receiver	gridpoint *Gridpoint 网格点结构体对象指针
// @return	uint64	网格点ID
func (gridpoint *Gridpoint) GetId() uint64 {
	t := &Gridpoint{}
	DB.Mysql.Where("coordinate_x = ? AND coordinate_y = ? AND coordinate_z = ?",
		gridpoint.Coordinate_x, gridpoint.Coordinate_y, gridpoint.Coordinate_z).Find(&t)
	return t.Id
}

// @title	GetGridpoint
// @description	根据场所id，x，y，z坐标查询网格点
// @auth	高宏宇
// @param	coordinate_x float64 网格点x坐标    coordinate_y float64 网格点y坐标    coordinate_z float64 网格点z坐标    place_id uint64 网格点ID
// @return	*Gridpoint 网格点对象指针    error 错误信息
func GetGridpoint(coordinate_x float64, coordinate_y float64, coordinate_z float64, place_id uint64) (*Gridpoint, error) {
	gridpoint := &Gridpoint{}
	// 查询结果为空报错record not found
	db := DB.Mysql.Where("coordinate_x = ? AND coordinate_y = ? AND coordinate_z = ? AND place_id = ?",
		coordinate_x, coordinate_y, coordinate_z, place_id).Find(&gridpoint)
	return gridpoint, db.Error
}

// @title	GetGridpointById
// @description	根据ID查询网格点
// @auth	高宏宇
// @param	id uint64 网格点ID
// @return	*Gridpoint 网格点对象指针    error 错误信息
func GetGridpointById(id uint64) (*Gridpoint, error) {
	t := &Gridpoint{}
	db := DB.Mysql.Where("id = ?", id).Find(&t)
	return t, db.Error
}
