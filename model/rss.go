package model

import (
	"time"
)

type Rss struct {
	Id                 uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Rss                int       `json:"rss" gorm:"column:rss;not null" binding:"required"`
	Reference_point_id uint64    `json:"reference_point_id"`
	Ap_id              uint64    `json:"ap_id"`
	Createdate         time.Time `gorm:"column:createdate"`
	Updatedate         time.Time `gorm:"column:updatedate"`
}

// 向数据库插入场所
func (rss *Rss) Create() error {
	return DB.Mysql.Create(&rss).Error
}

func (rss *Rss) GetId() uint64 {
	t := &Referencepoint{}
	// TODO 添加查询失败时的处理
	DB.Mysql.Where("reference_point_id = ? AND ap_id = ?",
		rss.Reference_point_id, rss.Ap_id).Find(&t)
	return t.Id
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
