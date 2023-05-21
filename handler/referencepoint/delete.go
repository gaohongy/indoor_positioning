package referencepoint

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// @title	Delete
// @description	删除参考点API
// @auth	高宏宇
// @param	ctx *gin.Context
func Delete(ctx *gin.Context) {
	log.Info("Referencepoint Delete function called")

	// 解析请求参数
	id_str := ctx.Query("id")         // id_str: 字符串类型
	id_int, _ := strconv.Atoi(id_str) // id_int: int类型
	id_uint64 := uint64(id_int)       // id_uint64: uint64类型

	// 删除参考点
	err := model.DeleteReferencepoint(id_uint64)
	// 查询结果为空err.Error() = "record not found"
	if err != nil {
		log.Error("referencepoint delete error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// 发送响应
	handler.SendResponse(ctx, nil, nil)
}
