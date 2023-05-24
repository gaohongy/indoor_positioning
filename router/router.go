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

	//	接口权限声明：
	//	public：无需token认证
	//	private_all：需token认证（即需登录用户），管理员和普通用户均可访问
	//	private_admin：需token认证，仅管理员可访问
	g.POST("/user", user.Create)   // 用户注册（public）
	g.POST("/session", user.Login) // 用户登录（public）

	// 需token认证，管理员和普通用户均可访问的用户相关API
	u_pirvate_all := g.Group("/user")
	u_pirvate_all.Use(middleware.GeneralAuthMiddleware())
	{

		u_pirvate_all.PUT("/place_id", user.PutPlaceId) // 登录用户修改自身场所id（private_all）
		u_pirvate_all.GET("/info", user.GetInfo)        // 登录用户获取自身信息（private_all）
	}

	// 需token认证，仅管理员可访问的用户相关API
	u_pirvate_admin := g.Group("/user")
	u_pirvate_admin.Use(middleware.AdminAuthMiddleware())
	{
		u_pirvate_admin.PUT("", user.Put)                  // 管理员修改同一场所用户的信息（private_admin）
		u_pirvate_admin.GET("", user.Get)                  // 用户管理界面获取当前场所全部用户信息（pirvate_admin）
		u_pirvate_admin.DELETE("", user.Delete)            // 删除用户（private_admin）
		u_pirvate_admin.GET("/count", user.GetCount)       // 获取当前场所不同时间点各类用户数量（private_admin）
		u_pirvate_admin.GET("/location", user.GetLocation) // 获取用户最新的位置信息（private_admin）
	}

	// 无需认证的用户相关API
	g.GET("/place", place.Get) // 获取用户列表（public）

	// 需token认证，管理员和普通用户均可访问的场所相关API
	p_private_all := g.Group("/place")
	p_private_all.Use(middleware.GeneralAuthMiddleware())
	{
		p_private_all.POST("/gridpoint", gridpoint.Create) // 添加网格点（private_all)
		p_private_all.POST("/pathpoint", pathpoint.Create) // 添加路径点（private_all)
		p_private_all.GET("/map", place.GetMap)            // 查询场所对应地图ID（private_all）

	}

	// 需token认证，仅管理员可访问的场所相关API
	p_private_admin := g.Group("/place")
	p_private_admin.Use(middleware.AdminAuthMiddleware())
	{
		p_private_admin.POST("", place.Create)                           // 添加场所（private_admin）
		p_private_admin.POST("/ap", ap.Create)                           // 添加AP（private_admin)
		p_private_admin.GET("/ap", ap.Get)                               // 获取符合时间条件的AP列表（private_admin)
		p_private_admin.PUT("/ap", ap.Put)                               // 修改AP信息（private_admin）
		p_private_admin.DELETE("/ap", ap.Delete)                         // 删除AP（private_admin）
		p_private_admin.POST("/referencepoint", referencepoint.Create)   // 创建参考点（private_admin）
		p_private_admin.GET("/referencepoint", referencepoint.Get)       // 获取符合时间条件的参考点列表（private_admin）
		p_private_admin.DELETE("/referencepoint", referencepoint.Delete) // 删除参考点（private_admin）
		p_private_admin.PUT("/referencepoint", referencepoint.Put)       // 修改参考点信息（private_admin）
		p_private_admin.GET("/pathpoint", pathpoint.Get)                 // 获取不同用户的路径点列表（private_admin）
	}

	// 需token认证，管理员和普通用户均可访问的定位API
	l := g.Group("/location")
	l.Use(middleware.GeneralAuthMiddleware())
	{
		l.GET("", location.Get) // 获取定位信息（private_all）
	}

	return g
}
