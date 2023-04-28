package middleware

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
)

// 目前能够正常请求的条件有二：1.请求header中Authorization字段非空 2.Authorization确实是经过jwt_secret签发得到的token
// private_all权限
func GeneralAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(ctx); err != nil {
			handler.SendResponse(ctx, errno.ErrorTokenInvalid, nil)
			ctx.Abort() // 终止当前中间件以后的中间件执行，但是会执行当前中间件的后续逻辑
			return
		}

		ctx.Next()
	}
}

// private_admin权限
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse the json web token.
		var content *token.Context
		var err error

		if content, err = token.ParseRequest(ctx); err != nil {
			handler.SendResponse(ctx, errno.ErrorTokenInvalid, nil)
			ctx.Abort() // 终止当前中间件以后的中间件执行，但是会执行当前中间件的后续逻辑
			return
		}

		user, _ := model.GetUserById(content.ID)
		if user.Usertype != 0 {
			handler.SendResponse(ctx, errno.ErrorTokenInvalid, nil)
			ctx.Abort() // 终止当前中间件以后的中间件执行，但是会执行当前中间件的后续逻辑
			return
		}

		ctx.Next()
	}
}
