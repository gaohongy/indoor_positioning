package user

import (
	"indoor_positioning/model"
	"time"
)

// user相关api所需请求响应结构

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType int    `json:"usertype"`
}

// TODO 确定返回数据
type CreateResponse struct {
	Username string `json:"username"`
}

type PutPlaceIdRequest struct {
	Place_id uint64 `json:"place_id"`
}

type PutRequest struct {
	Id       uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Usertype int    `json:"usertype" validate:"required"`
}

type GetResponse struct {
	User_list []model.User_Brief `json:"user_list"`
}

type LoginResponse struct {
	UserType int    `json:"usertype"`
	Place_id uint64 `json:"place_id"`
	Token    string `json:"token"`
}

type GetUserAmountResponse struct {
	Date               time.Time `json:"date"`
	AdminAmount        uint64    `json:"adminAmount"`
	OrdinaryUserAmount uint64    `json:"ordinaryUserAmount"`
	SumUserAmount      uint64    `json:"sumUserAmount"`
}

type UserAmount struct {
	AdminAmount        uint64 `json:"adminAmount"`
	OrdinaryUserAmount uint64 `json:"ordinaryUserAmount"`
	SumUserAmount      uint64 `json:"sumUserAmount"`
}
