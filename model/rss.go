package model

import (
	"time"

	"github.com/zxmrlc/log"
)

type Rss struct {
	Id                 uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Rss                float64   `json:"rss" gorm:"column:rss;not null" binding:"required"`
	Reference_point_id uint64    `json:"reference_point_id" gorm:"column:reference_point_id"`
	Ap_id              uint64    `json:"ap_id" gorm:"column:ap_id"`
	Createdate         time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate         time.Time `json:"updatedate" gorm:"column:updatedate"`
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

// 根据参考点id获取所有rss
func ListRssByReferencepointid(reference_point_id uint64, limit, offset int) ([]*Rss, uint64, error) {
	if limit == 0 {
		limit = 50
	}

	// TODO 这里list的大小为0？
	rss_list := make([]*Rss, 0)
	var count uint64

	// TODO 目前暂时不考虑分页,注释代码为添加分页的版本
	// TODO "error":"Error 1054: Unknown column 'reference_point_id' in 'where clause'" referencepoint.go中之前同样有这个错误
	// 这里是统计查询结果行数，错误是因为Model中传递的类型不对，它会根据Model中指定的类型去对应数据库表
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

	// TODO 这里最初用reference_point_id就查不到结果，当替换为数字后就能查到，再换为reference_point_id就可以查到了，不懂为什么
	if err := DB.Mysql.Where("reference_point_id = ?", reference_point_id).Find(&rss_list).Error; err != nil {
		return rss_list, count, err
	}
	// if err := DB.Mysql.Where("reference_point_id = ?", reference_point_id).Order("id desc").Find(&rss_list).Error; err != nil {
	// 	return rss_list, count, err
	// }

	return rss_list, count, nil
}
