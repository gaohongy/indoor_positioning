package model

// 和user相关的数据库的接口函数
import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TODO validate需要细化
// TODO 时间充裕考虑自行实现生成salt以及加密过程
// User represents a registered user.
// 标签中的validate即给出该参数的校验规则
// 添加db.SingularTable(true)后，user即为对应数据库表名
type User struct {
	Id         uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	Username   string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password   string    `gorm:"column:pwdhash" json:"pwdhash" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
	Usertype   int       `json:"usertype" binding:"required" validate:"required"`
	Place_id   int       `json:"-"`
	Createdate time.Time `gorm:"column:createdate"`
	Updatedate time.Time `gorm:"column:updatedate"`
}

// 向数据库插入用户
func (user *User) Create() error {
	return DB.Mysql.Create(&user).Error
}

// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (user *User) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }

// Encrypt the user password.
func (user *User) Encrypt() (err error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedBytes)
	return err
}

// // Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
// func (u *User) Compare(pwd string) (err error) {
// 	err = auth.Compare(u.Password, pwd)
// 	return
// }
