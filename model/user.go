package model

// 和user相关的数据库的接口函数
import (
	"indoor_positioning/pkg/auth"
	"time"
)

// TODO validate需要细化
// TODO 时间充裕考虑自行实现生成salt以及加密过程
// TODO usertype在注册时需要绑定，但在登录时候无需绑定，是否考虑手动校验
// User represents a registered user.
// 标签中的validate即给出该参数的校验规则
// 添加db.SingularTable(true)后，user即为对应数据库表名
type User struct {
	Id         uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Username   string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password   string    `json:"password" gorm:"column:pwdhash;not null" binding:"required" validate:"min=5,max=128"`
	Usertype   int       `json:"usertype" validate:"required"`
	Place_id   uint64    `json:"place_id"`
	Createdate time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate time.Time `json:"updatedate" gorm:"column:updatedate"`
}

type User_Brief struct {
	Id         uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Username   string    `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Usertype   int       `json:"usertype" validate:"required"`
	Place_id   uint64    `json:"place_id"`
	Createdate time.Time `json:"createdate" gorm:"column:createdate"`
	Updatedate time.Time `json:"updatedate" gorm:"column:updatedate"`
}

// TODO User的方法必须要在同一文件中生成，所以Encrypt和Compare必须写在这里，为了增强代码易读性，将具体实现放置于pkg的auth下
// Encrypt the user password.
func (user *User) Encrypt() (err error) {
	user.Password, err = auth.Encrypt(user.Password)
	return err
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (user *User) Compare(password string) (err error) {
	err = auth.Compare(user.Password, password)
	return err
}

// -------------------------------------------------------------------------------------------------

// 向数据库插入用户
func (user *User) Create() error {
	return DB.Mysql.Create(&user).Error
}

// 根据id获取user
func GetUserById(id uint64) (*User, error) {
	user := &User{}
	db := DB.Mysql.Where("id = ?", id).Find(&user) // 查询结果值存储在user中，返回值是个*gorm.DB
	return user, db.Error
}

// 根据username获取user
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	// TODO 这里的返回值不太能理解为什么是*gorm.DB
	db := DB.Mysql.Where("username = ?", username).Find(&user) // 查询结果值存储在user中，返回值是个*gorm.DB
	return user, db.Error
}

// 修改用户place_id
func (user *User) Update(place_id uint64) error {
	db := DB.Mysql.Model(user).Update("place_id", place_id)
	return db.Error
}

func GetUserByPlaceId(place_id int) ([]*User, error) {
	user_list := make([]*User, 0)
	db := DB.Mysql.Where("place_id = ?", place_id).Find(&user_list)
	return user_list, db.Error
}

func FilterUserByTime(place_id int, begin_time time.Time, end_time time.Time) ([]*User, error) {
	user_list := make([]*User, 0)

	db := DB.Mysql.Where("place_id = ? AND createdate BETWEEN ? AND ?", place_id, begin_time, end_time).Find(&user_list)

	return user_list, db.Error
}

// TODO 应当验证用户名是否重复
// 结构体属性合法性校验
// 目前仅校验Username,Password,Usertype
// func (user *User) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(user)
// }
