package model

import (
	"time"

	"github.com/zxmrlc/log"
)

// 参考点信息
type Referencepoint struct {
	Id            uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"` // 参考点ID
	Grid_point_id uint64    `json:"grid_point_id" gorm:"column:grid_point_id"`      // 网格点ID
	Place_id      uint64    `json:"place_id" gorm:"column:place_id"`                // 场所ID
	Createdate    time.Time `json:"createdate" gorm:"column:createdate"`            // 创建时间
	Updatedate    time.Time `json:"updatedate" gorm:"column:updatedate"`            // 修改时间
}

// 参考点信息（包含坐标详细信息）
type Referencepoint_Detail struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`                      // 参考点ID
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"` // 参考点x坐标
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"` // 参考点y坐标
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"` // 参考点z坐标
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`                                 // 创建时间
	Updatedate   time.Time `json:"updatedate" gorm:"column:updatedate"`                                 // 修改时间
}

// @title	Create
// @description	结构体方法，向数据库插入参考点
// @auth	高宏宇
// @receiver	referencepoint *Referencepoint 参考点结构体对象指针
// @return	error 错误信息
func (referencepoint *Referencepoint) Create() error {
	return DB.Mysql.Create(&referencepoint).Error
}

// @title	GetId
// @description	结构体方法，查询参考点ID
// @auth	高宏宇
// @receiver	referencepoint *Referencepoint 参考点结构体对象指针
// @return	uint64	参考点ID
func (referencepoint *Referencepoint) GetId() uint64 {
	t := &Referencepoint{}
	DB.Mysql.Where("place_id = ? AND grid_point_id = ?",
		referencepoint.Place_id, referencepoint.Grid_point_id).Find(&t)
	return t.Id
}

// @title	Update
// @description	结构体方法，修改参考点信息
// @auth	高宏宇
// @receiver	referencepoint *Referencepoint 参考点结构体对象指针
// @param	grid_point_id uint64 网格点ID
// @return	error	错误信息
func (referencepoint *Referencepoint) Update(grid_point_id uint64) error {
	db := DB.Mysql.Model(referencepoint).Update(map[string]interface{}{"grid_point_id": grid_point_id})
	return db.Error
}

// @title	GetReferencepointById
// @description	根据ID查询参考点
// @auth	高宏宇
// @param	id uint64 参考点ID
// @return	*Referencepoint 参考点对象指针    error 错误信息
func GetReferencepointById(id uint64) (*Referencepoint, error) {
	referencepoint := &Referencepoint{}
	db := DB.Mysql.Where("id = ?", id).Find(&referencepoint)
	return referencepoint, db.Error
}

// @title	ListReferencepointByPlaceid
// @description	查询全部参考点
// @auth	高宏宇
// @param	place_id uint64 场所ID    limit int 每页的结果条目数量    offset int 从第几条记录开始(mysql中第1条索引为0)
// @return	[]*Referencepoint 参考点对象指针数组    uint64 查询结果条目数    error 错误信息
func ListReferencepointByPlaceid(place_id uint64, limit, offset int) ([]*Referencepoint, uint64, error) {
	if limit == 0 {
		limit = 50
	}

	referencepoints := make([]*Referencepoint, 0)
	var count uint64

	// 暂时不考虑分页,注释代码为添加分页的版本
	if err := DB.Mysql.Model(&Referencepoint{}).Where("place_id = ?", place_id).Count(&count).Error; err != nil {
		log.Error("reference count error", err)
		return referencepoints, count, err
	}

	// if err := DB.Mysql.Where("place_id = ?", place_id).Offset(offset).Limit(limit).Order("id desc").Find(&referencepoints).Error; err != nil {
	// 	return referencepoints, count, err
	// }

	if err := DB.Mysql.Where("place_id = ?", place_id).Find(&referencepoints).Error; err != nil {
		return referencepoints, count, err
	}
	// if err := DB.Mysql.Where("place_id = ?", place_id).Order("id desc").Find(&referencepoints).Error; err != nil {
	// 	return referencepoints, count, err
	// }

	return referencepoints, count, nil
}

// @title	FilterReferencepointByTime
// @description	筛选指定场所中指定时间段内的参考点
// @auth	高宏宇
// @param	place_id int 场所ID    begin_time time.Time 筛选开始时间    end_time time.Time 筛选结束时间
// @return	[]*Referencepoint 参考点对象指针数组    uint64 查询结果条目数    error 错误信息
func FilterReferencepointByTime(place_id int, begin_time time.Time, end_time time.Time) ([]*Referencepoint, error) {
	referencepoint_list := make([]*Referencepoint, 0)
	db := DB.Mysql.Where("place_id = ? AND createdate BETWEEN ? AND ?", place_id, begin_time, end_time).Find(&referencepoint_list)
	return referencepoint_list, db.Error
}

// @title	DeleteReferencepoint
// @description	删除指定参考点
// @auth	高宏宇
// @param	id uint64 参考点ID
// @return	error 错误信息
func DeleteReferencepoint(id uint64) error {
	referencepoint := &Referencepoint{}
	db := DB.Mysql.Where("id = ?", id).Delete(&referencepoint)
	return db.Error
}
