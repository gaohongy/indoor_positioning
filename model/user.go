package model

// 和user相关的数据库的接口函数
import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

// TODO validate需要细化
// User represents a registered user.
// 标签中的validate即给出该参数的校验规则
type UserModel struct {
	Id        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Username  string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Salt      string    `json:"-"`
	Password  string    `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	Usertype  int       `json:"usertype" binding:"required" validate:"required"`
	Place_id  int       `json:"-"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"-"`
}

// 向数据库插入用户
func (user *UserModel) Create() error {
	return DB.Mysql.Create(&user).Error
}

// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
func (user *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}

// Encrypt the user password.
// func (u *UserModel) Encrypt() (err error) {
// 	u.Password, err = auth.Encrypt(u.Password)
// 	return
// }

// // Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
// func (u *UserModel) Compare(pwd string) (err error) {
// 	err = auth.Compare(u.Password, pwd)
// 	return
// }
