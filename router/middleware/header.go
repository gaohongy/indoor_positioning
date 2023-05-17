package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @title	NoCache
// @description	NoCache is a middleware function that appends headers to prevent the client from caching the HTTP response.
// @auth	高宏宇
// @param	c	*gin.Context	Context指针
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-store")                               // 不缓存
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")                // 服务器告知浏览器资源失效时间，超过该时间后，资源失效
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat)) // 服务器告知浏览器资源最后修改时间
	c.Next()
}

// @title Options
// @description	处理options请求中间件函数(preflight request 预检请求)，设置请求头中相关权限字段
// @auth	高宏宇
// @param	c	*gin.Context	Context指针
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		// 设置响应标头
		c.Header("Access-Control-Allow-Origin", "*")                                            // 发生跨域请求时，浏览器仅支持指定源发起的请求(*：任意源)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")           // 发生跨域请求时，浏览器仅支持指定的方法发起的请求
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept") // 发生跨域请求时，浏览器支持请求首部中包含的字段
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")                             // 服务器接受的请求方法
		c.Header("Content-Type", "application/json")                                            // 服务器返回的内容类型, application表明是某种二进制数据
		c.AbortWithStatus(200)
	}
}

// @title Secure
// @description	Secure is a middleware function that appends security and resource access headers.
// @auth	高宏宇
// @param	c	*gin.Context	Context指针
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")  // 发生跨域请求时，浏览器仅支持指定源发起的请求(*：任意源)
	c.Header("X-Frame-Options", "DENY")           //  禁止其他页面以 frame 形式嵌入
	c.Header("X-Content-Type-Options", "nosniff") // 提示客户端一定要遵循在 Content-Type 首部中对 MIME 类型 的设定
	c.Header("X-XSS-Protection", "1; mode=block") // 当检测到跨站脚本攻击 (XSS (en-US)) 时，浏览器将停止加载页面. 1;mode=block:启用 XSS 过滤。如果检测到攻击，浏览器将不会清除页面，而是阻止页面加载。
	if c.Request.TLS != nil {                     // TLS: Transport Layer Security 安全传输层协议，可看为SSL 3.0的新版本
		c.Header("Strict-Transport-Security", "max-age=31536000") // 通知浏览器应该只通过 HTTPS 访问该站点，并且以后使用 HTTP 访问该站点的所有尝试都应自动重定向到 HTTPS；max-age：浏览器记忆的只能使用 HTTPS 访问站点的最大时间量
	}
}
