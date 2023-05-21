package user

import (
	"indoor_positioning/model"
	"time"
)

// 创建用户API请求数据结构
type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType int    `json:"usertype"`
}

// 创建用户API响应数据结构
type CreateResponse struct {
	Username string `json:"username"`
}

// 登录用户修改所在场所API请求数据结构
type PutPlaceIdRequest struct {
	Place_id uint64 `json:"place_id"`
}

// 修改用户信息API请求数据结构
type PutRequest struct {
	Id       uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Usertype int    `json:"usertype" validate:"required"`
}

// 查询符合时间条件的用户API响应数据结构
type GetResponse struct {
	User_list []model.User_Brief `json:"user_list"`
}

// 用户登录API响应数据结构
type LoginResponse struct {
	UserType int    `json:"usertype"`
	Place_id uint64 `json:"place_id"`
	Token    string `json:"token"`
}

// 查询用户数量API响应数据结构
type GetUserAmountResponse struct {
	Date               time.Time `json:"date"`
	AdminAmount        uint64    `json:"adminAmount"`
	OrdinaryUserAmount uint64    `json:"ordinaryUserAmount"`
	SumUserAmount      uint64    `json:"sumUserAmount"`
}

// 用户数量信息
type UserAmount struct {
	AdminAmount        uint64 `json:"adminAmount"`        // 管理员数量
	OrdinaryUserAmount uint64 `json:"ordinaryUserAmount"` // 普通用户数量
	SumUserAmount      uint64 `json:"sumUserAmount"`      // 总人数
}
