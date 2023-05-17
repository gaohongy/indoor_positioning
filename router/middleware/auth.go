package middleware

import (
	"indoor_positioning/handler"
	"indoor_positioning/model"
	"indoor_positioning/pkg/errno"
	"indoor_positioning/pkg/token"

	"github.com/gin-gonic/gin"
)

// @title	GeneralAuthMiddleware
// @description	private_all权限认证中间件，供需token认证，管理员和普通用户均可访问的相关API使用。认证通过条件：1.请求header中Authorization字段非空 2.Authorization是经过jwt_secret签发得到的token
// @auth	高宏宇
// @return	gin.HandlerFunc
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

// @title	GeneralAuthMiddleware
// @description	private_admin权限认证中间件，供需token认证，仅管理员可访问的相关API使用。认证通过条件：1.请求header中Authorization字段非空 2.Authorization是经过jwt_secret签发得到的token 3.请求发起用户为管理员
// @auth	高宏宇
// @return	gin.HandlerFunc
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
		if user.Usertype != 0 { // 请求发起用户不是管理员
			handler.SendResponse(ctx, errno.ErrorTokenInvalid, nil)
			ctx.Abort() // 终止当前中间件以后的中间件执行，但是会执行当前中间件的后续逻辑
			return
		}

		ctx.Next()
	}
}
