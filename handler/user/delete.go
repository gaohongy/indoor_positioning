package user

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zxmrlc/log"
)

// Create creates a new user account.
func Delete(ctx *gin.Context) {
	log.Info("User Delete function called")

	// 解析请求body
	id_str := ctx.Query("id")
	id_int, _ := strconv.Atoi(id_str)
	id_uint64 := uint64(id_int)

	// 处理参考点
	err := model.DeleteUser(id_uint64)
	// 查询结果为空err.Error() = "record not found"
	if err != nil {
		log.Error("user delete error", err)
		handler.SendResponse(ctx, errno.ErrorDatabase, nil)
		return
	}

	// TODO 验证参数合法性
	// if err := user.Validate(); err != nil {
	// 	handler.SendResponse(ctx, errno.ErrorValidation, nil)
	// 	return
	// }

	// createResponse := CreateResponse{
	// 	Username: request.Username,
	// }

	// 发送响应
	handler.SendResponse(ctx, nil, nil)
}
