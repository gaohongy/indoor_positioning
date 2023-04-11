package router

import (
	"net/http"

	"indoor_positioning/handler/ap"
	"indoor_positioning/handler/gridpoint"
	"indoor_positioning/handler/location"
	"indoor_positioning/handler/pathpoint"
	"indoor_positioning/handler/place"
	"indoor_positioning/handler/referencepoint"
	"indoor_positioning/handler/user"
	"indoor_positioning/router/middleware"

	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// set middlewares
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	// set 404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	g.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.POST("/user", user.Create)   // 用户注册
	g.POST("/session", user.Login) // 用户登录

	u := g.Group("/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.PUT("", user.Put)
		u.GET("", user.Get)
	}

	// TODO 添加管理员身份认证中间件，但是这里的路由需要细化，因为普通用户是有添加路径点的权限的，那么自然要有添加网格点的权限
	p := g.Group("/place")
	p.Use(middleware.AuthMiddleware())
	{
		p.POST("", place.Create)
		p.POST("/ap", ap.Create)
		p.GET("/ap", ap.Get)
		p.POST("/referencepoint", referencepoint.Create)
		p.GET("/referencepoint", referencepoint.Get)
		p.POST("/gridpoint", gridpoint.Create)
		p.POST("/pathpoint", pathpoint.Create)
	}

	l := g.Group("/location")
	l.Use(middleware.AuthMiddleware())
	{
		l.GET("", location.Get)
	}

	return g
}
