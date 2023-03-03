package model

import (
	"time"

	"github.com/zxmrlc/log"
)

type Referencepoint struct {
	Id            uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Grid_point_id uint64    `json:"grid_point_id" gorm:"column:grid_point_id"`
	Place_id      uint64    `json:"place_id" gorm:"column:place_id"`
	Createdate    time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate    time.Time `json:"updatedate" gorm:"column:updatedate"`
}

// 向数据库插入场所
func (referencepoint *Referencepoint) Create() error {
	return DB.Mysql.Create(&referencepoint).Error
}

func (referencepoint *Referencepoint) GetId() uint64 {
	t := &Referencepoint{}
	// TODO 添加查询失败时的处理
	DB.Mysql.Where("place_id = ? AND grid_point_id = ?",
		referencepoint.Place_id, referencepoint.Grid_point_id).Find(&t)
	return t.Id
}

// TODO 经纬度和地址验证
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (place *Place) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }

// ListUser List all users
// TODO 这里的查询方式是考虑了分页的，在结果集中数量较多时提高查询速度
// limit：每页的结果条目数量 offset：从第几条记录开始(mysql中第1条索引为0)
func ListReferencepointByPlaceid(place_id uint64, limit, offset int) ([]*Referencepoint, uint64, error) {
	if limit == 0 {
		limit = 50
	}

	// TODO 这里大小定为0？
	referencepoints := make([]*Referencepoint, 0)
	var count uint64

	// TODO 目前暂时不考虑分页,注释代码为添加分页的版本
	if err := DB.Mysql.Model(&Referencepoint{}).Where("place_id = ?", place_id).Count(&count).Error; err != nil {
		log.Error("reference count error", err)
		return referencepoints, count, err
	}
	// if err := DB.Mysql.Where("place_id = ?", place_id).Offset(offset).Limit(limit).Order("id desc").Find(&referencepoints).Error; err != nil {
	// 	return referencepoints, count, err
	// }

	// TODO 这里最初用reference_point_id就查不到结果，当替换为数字后就能查到，再换为reference_point_id就可以查到了，不懂为什么
	if err := DB.Mysql.Where("place_id = ?", place_id).Find(&referencepoints).Error; err != nil {
		return referencepoints, count, err
	}
	// if err := DB.Mysql.Where("place_id = ?", place_id).Order("id desc").Find(&referencepoints).Error; err != nil {
	// 	return referencepoints, count, err
	// }

	return referencepoints, count, nil
}
