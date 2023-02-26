package router

import (
	"net/http"

	"indoor_positioning/handler/ap"
	"indoor_positioning/handler/gridpoint"
	"indoor_positioning/handler/location"
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

	p := g.Group("/place")
	// p.Use(middleware.AuthMiddleware())
	{
		p.POST("", place.Create)
		p.POST("/ap", ap.Create)
		p.POST("/referencepoint", referencepoint.Create)
		p.POST("/gridpoint", gridpoint.Create)
	}

	l := g.Group("/location")
	// l.Use(middleware.AuthMiddleware())
	{
		l.GET("", location.Get)
	}

	return g
}
