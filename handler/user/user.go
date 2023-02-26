package user

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
