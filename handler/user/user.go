package user

import "indoor_positioning/model"

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

type PutRequest struct {
	Place_id uint64 `json:"place_id"`
}

type GetResponse struct {
	User_list []model.User_Brief `json:"user_list"`
}

type LoginResponse struct {
	UserType int    `json:"usertype"`
	Place_id uint64 `json:"place_id"`
	Token    string `json:"token"`
}
