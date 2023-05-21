package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 路径点信息
type Pathpoint struct {
	Id            uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"` // 路径点ID
	User_id       uint64    `json:"user_id"`                                        // 用户ID
	Grid_point_id uint64    `json:"grid_point_id"`                                  // 网格点ID
	Place_id      uint64    `json:"place_id"`                                       // 场所ID
	Createdate    time.Time `json:"createdate" gorm:"column:createdate"`            // 创建时间
	Updatedate    time.Time `json:"updatedate" gorm:"column:updatedate"`            // 修改时间
}

// 路径点信息（包含坐标详细信息）
type Pathpoint_Detail struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`                      // 路径点ID
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"` // 路径点x坐标
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"` // 路径点y坐标
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"` // 路径点z坐标
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`                                 // 创建时间
}

// @title	Create
// @description	结构体方法，向数据库插入路径点
// @auth	高宏宇
// @receiver	pathpoint *Pathpoint 路径点结构体对象指针
// @return	error	错误信息
func (pathpoint *Pathpoint) Create() error {
	return DB.Mysql.Create(&pathpoint).Error
}

// @title	FilterPathpointByTimeAndUser
// @description	筛选指定场所中指定时间段内指定用户的路径点
// @auth	高宏宇
// @param	place_id int 场所ID    user_id int 用户ID    begin_time time.Time 筛选开始时间    end_time time.Time 筛选结束时间
// @return	*[]Pathpoint 路径点对象数组指针    error 错误信息
func FilterPathpointByTimeAndUser(place_id int, user_id int, begin_time time.Time, end_time time.Time) (*[]Pathpoint, error) {
	pathpoint_list := &[]Pathpoint{}
	var db *gorm.DB

	// 无时间筛选条件
	if begin_time.IsZero() && end_time.IsZero() {
		db = DB.Mysql.Where("place_id = ? AND user_id = ?", place_id, user_id).Find(&pathpoint_list)

		// 有时间筛选条件
	} else {
		db = DB.Mysql.Where("place_id = ? AND user_id = ? AND createdate BETWEEN ? AND ?", place_id, user_id, begin_time, end_time).Find(&pathpoint_list)
	}

	return pathpoint_list, db.Error
}

// @title	FilterLatestPathpointByUserId
// @description	查询用户最新路径点
// @auth	高宏宇
// @param	user_id uint64 用户ID
// @return	*Pathpoint 路径点对象指针    error 错误信息
func FilterLatestPathpointByUserId(user_id uint64) (*Pathpoint, error) {
	pathpoint := &Pathpoint{}
	db := DB.Mysql.Where("user_id = ?", user_id).Order("createdate").Find(&pathpoint)
	return pathpoint, db.Error
}
