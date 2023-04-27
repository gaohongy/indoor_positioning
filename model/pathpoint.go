package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Pathpoint struct {
	Id            uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	User_id       uint64    `json:"user_id"`
	Grid_point_id uint64    `json:"grid_point_id"`
	Place_id      uint64    `json:"place_id"`
	Createdate    time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate    time.Time `json:"updatedate" gorm:"column:updatedate"`
}

type Pathpoint_Detail struct {
	Id           uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Coordinate_x float64   `json:"coordinate_x" gorm:"column:coordinate_x;not null" binding:"required"`
	Coordinate_y float64   `json:"coordinate_y" gorm:"column:coordinate_y;not null" binding:"required"`
	Coordinate_z float64   `json:"coordinate_z" gorm:"column:coordinate_z;not null" binding:"required"`
	Createdate   time.Time `json:"createdate" gorm:"column:createdate"`
}

// 向数据库插入路径点
func (pathpoint *Pathpoint) Create() error {
	return DB.Mysql.Create(&pathpoint).Error
}

func FilterPathpointByTimeAndUser(place_id int, user_id int, begin_time time.Time, end_time time.Time) (*[]Pathpoint, error) {
	pathpoint_list := &[]Pathpoint{}
	var db *gorm.DB

	if begin_time.IsZero() && end_time.IsZero() { // 无时间筛选条件
		db = DB.Mysql.Where("place_id = ? AND user_id = ?", place_id, user_id).Find(&pathpoint_list)
	} else { // 有时间筛选条件
		db = DB.Mysql.Where("place_id = ? AND user_id = ? AND createdate BETWEEN ? AND ?", place_id, user_id, begin_time, end_time).Find(&pathpoint_list)
	}

	return pathpoint_list, db.Error
}

func FilterLatestPathpointByUserId(user_id uint64) (*Pathpoint, error) {
	pathpoint := &Pathpoint{}
	// TODO find是查找第一条数据，这里要查询的是最近一次的路径点，应当是时间戳更大的，因此按照createdate降序desc排列拿到的第一条数据才是需要的，但是这里结果是反的，尚不清楚原因
	db := DB.Mysql.Where("user_id = ?", user_id).Order("createdate").Find(&pathpoint)
	return pathpoint, db.Error
}

// func (pathpoint *Pathpoint) GetId() uint64 {
// 	t := &Referencepoint{}
// 	// TODO 添加查询失败时的处理
// 	DB.Mysql.Where("place_id = ? AND grid_point_id = ?",
// 		referencepoint.Place_id, referencepoint.Grid_point_id).Find(&t)
// 	return t.Id
// }

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
