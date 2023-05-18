package model

// 和user相关的数据库的接口函数
import (
	"indoor_positioning/pkg/auth"
	"time"
)

// 用户信息
type User struct {
	Id         uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`                                      // 用户ID
	Username   string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"` // 用户名
	Password   string    `json:"password" gorm:"column:pwdhash;not null" binding:"required" validate:"min=5,max=128"` // 用户密码
	Usertype   int       `json:"usertype" validate:"required"`                                                        // 用户类型
	Place_id   uint64    `json:"place_id"`                                                                            // 用户所在场所ID
	Createdate time.Time `json:"createdate" gorm:"column:createdate"`                                                 // 创建时间
	Updatedate time.Time `json:"updatedate" gorm:"column:updatedate"`                                                 // 修改时间
}

// 用户信息（简要信息）
type User_Brief struct {
	Id         uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`                                      // 用户ID
	Username   string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"` // 用户名
	Usertype   int       `json:"usertype" validate:"required"`                                                        // 用户类型
	Place_id   uint64    `json:"place_id"`                                                                            // 用户所在场所ID
	Createdate time.Time `json:"createdate" gorm:"column:createdate"`                                                 // 创建时间
	Updatedate time.Time `json:"updatedate" gorm:"column:updatedate"`                                                 // 修改时间
}

// @title	Encrypt
// @description	结构体方法，加密用户密码
// @auth	高宏宇
// @receiver	user *User 用户结构体对象指针
// @return	err error 错误信息
func (user *User) Encrypt() (err error) {
	user.Password, err = auth.Encrypt(user.Password)
	return err
}

// @title	Create
// @description	结构体方法，向数据库插入用户
// @auth	高宏宇
// @receiver	user *User 用户结构体对象指针
// @return	err error 错误信息
func (user *User) Create() error {
	return DB.Mysql.Create(&user).Error
}

// @title	UpdateUserPlaceId
// @description	结构体方法，修改用户所在场所
// @auth	高宏宇
// @receiver	user *User 用户结构体对象指针
// @param    place_id uint64 场所ID
// @return	err error 错误信息
func (user *User) UpdateUserPlaceId(place_id uint64) error {
	db := DB.Mysql.Model(user).Update("place_id", place_id)
	return db.Error
}

// @title	UpdateUserInfo
// @description	结构体方法，修改用户信息
// @auth	高宏宇
// @receiver	user *User 用户结构体对象指针
// @param    username string 用户名    usertype int 用户类型
// @return	err error 错误信息
func (user *User) UpdateUserInfo(username string, usertype int) error {
	db := DB.Mysql.Model(user).Update(map[string]interface{}{"username": username, "usertype": usertype})
	return db.Error
}

// @title	GetUserById
// @description	根据ID查询用户
// @auth	高宏宇
// @param	id uint64 用户ID
// @return	*User 用户对象指针    error 错误信息
func GetUserById(id uint64) (*User, error) {
	user := &User{}
	db := DB.Mysql.Where("id = ?", id).Find(&user) // 查询结果值存储在user中，返回值是个*gorm.DB
	return user, db.Error
}

// @title	GetUserByUsername
// @description	根据用户名查询用户
// @auth	高宏宇
// @param	username string 用户名
// @return	*User 用户对象指针    error 错误信息
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	db := DB.Mysql.Where("username = ?", username).Find(&user) // 查询结果值存储在user中，返回值是个*gorm.DB
	return user, db.Error
}

// @title	GetUserByPlaceId
// @description	查询场所全部用户
// @auth	高宏宇
// @param	place_id int 待查询场所ID
// @return	[]*User 用户对象指针数组    error 错误信息
func GetUserByPlaceId(place_id int) ([]*User, error) {
	user_list := make([]*User, 0)
	db := DB.Mysql.Where("place_id = ?", place_id).Find(&user_list)
	return user_list, db.Error
}

// @title	FilterUserByTime
// @description	筛选指定场所中指定时间段内的用户
// @auth	高宏宇
// @param	place_id int 待查询场所ID    begin_time time.Time 筛选开始时间    end_time time.Time 筛选结束时间
// @return	[]*User 用户对象指针数组    error 错误信息
func FilterUserByTime(place_id int, begin_time time.Time, end_time time.Time) ([]*User, error) {
	user_list := make([]*User, 0)

	db := DB.Mysql.Where("place_id = ? AND createdate BETWEEN ? AND ?", place_id, begin_time, end_time).Find(&user_list)

	return user_list, db.Error
}

// @title	DeleteUser
// @description	删除用户
// @auth	高宏宇
// @param	id uint64 待删除用户ID
// @return	error 错误信息
func DeleteUser(id uint64) error {
	user := &User{}
	db := DB.Mysql.Where("id = ?", id).Delete(&user)
	return db.Error
}
