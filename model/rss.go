package model

import (
	"time"

	"github.com/zxmrlc/log"
)

// WiFi指纹信息
type Rss struct {
	Id                 uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`      // WiFi指纹ID
	Rss                float64   `json:"rss" gorm:"column:rss;not null" binding:"required"`   // 接入点信号强度
	Reference_point_id uint64    `json:"reference_point_id" gorm:"column:reference_point_id"` // 参考点ID
	Ap_id              uint64    `json:"ap_id" gorm:"column:ap_id"`                           // 接入点ID
	Createdate         time.Time `json:"createdate" gorm:"column:createdate"`                 // 创建时间
	Updatedate         time.Time `json:"updatedate" gorm:"column:updatedate"`                 // 修改时间
}

// @title	Create
// @description	结构体方法，向数据库插入WiFi指纹
// @auth	高宏宇
// @receiver	rss *Rss WiFi指纹结构体对象指针
// @return	error	错误信息
func (rss *Rss) Create() error {
	return DB.Mysql.Create(&rss).Error
}

// @title	ListRssByReferencepointid
// @description	查询参考点全部WiFi指纹
// @auth	高宏宇
// @param	reference_point_id uint64 参考点ID    limit int 每页的结果条目数量    offset int 从第几条记录开始(mysql中第1条索引为0)
// @return	[]*Rss WiFi指纹对象指针数组    uint64 查询结果条目数    error 错误信息
func ListRssByReferencepointid(reference_point_id uint64, limit, offset int) ([]*Rss, uint64, error) {
	if limit == 0 {
		limit = 50
	}

	rss_list := make([]*Rss, 0)
	var count uint64

	// 注释代码为添加分页的版本
	if err := DB.Mysql.Model(&Rss{}).Where("reference_point_id = ?", reference_point_id).Count(&count).Error; err != nil {
		log.Error("error", err)
		return rss_list, count, err
	}
	// 等价版本
	// if err := DB.Mysql.Where("reference_point_id = ?", reference_point_id).Find(&rss_list).Count(&count).Error; err != nil {
	// 	log.Error("error", err)
	// 	return rss_list, count, err
	// }

	// if err := DB.Mysql.Where("place_id = ?", place_id).Offset(offset).Limit(limit).Order("id desc").Find(&referencepoints).Error; err != nil {
	// 	return referencepoints, count, err
	// }

	if err := DB.Mysql.Where("reference_point_id = ?", reference_point_id).Find(&rss_list).Error; err != nil {
		return rss_list, count, err
	}
	// if err := DB.Mysql.Where("reference_point_id = ?", reference_point_id).Order("id desc").Find(&rss_list).Error; err != nil {
	// 	return rss_list, count, err
	// }

	return rss_list, count, nil
}
