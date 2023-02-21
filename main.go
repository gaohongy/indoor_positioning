package main

import (
	"indoor_positioning/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	// 加载路由&设置中间件
	router.Load(
		// cores
		g,

		// set middlwares
	)

	http.ListenAndServe(":80", g)
}
