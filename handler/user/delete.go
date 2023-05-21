package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Delete
// @description	删除用户API
// @auth	高宏宇
// @param	ctx *gin.Context
func Delete(ctx *gin.Context) {
	log.Info("User Delete function called")

	// 解析请求参数
	id_str := ctx.Query("id")         // id_str: 字符串类型
	id_int, _ := strconv.Atoi(id_str) // id_int: int类型
	id_uint64 := uint64(id_int)       // id_uint64: uint64类型

	// 删除用户
	err := model.DeleteUser(id_uint64)
	// 删除失败
	if err != nil {
		log.Error("user delete error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 发送响应
	handler.SendResponse(ctx, nil, nil)
}
