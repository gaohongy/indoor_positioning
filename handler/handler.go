package handler

import (
	"net/http"

	"indoor_positioning/pkg/errno"

	"github.com/gin-gonic/gin"
)

// 通用响应模型
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	// 在DecodeErr()中实现了err为nil时解析为OK，所以err传nil时会返回ok
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
