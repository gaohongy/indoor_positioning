package main

import (
	"indoor_positioning/config"
	"indoor_positioning/router"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zxmrlc/log"
)

// 设置命令行参数
// 假设编译后得可执行文件为apiServer
// 则参数使用方式为 ./apiServer --config/-c "<file path>", 当不加选项时默认值为第3个参数空
var (
	cfg = pflag.StringP("config", "c", "", "config file path.")
)

func main() {
	// 打印默认选项
	pflag.PrintDefaults()
	// 解析命令行参数
	pflag.Parse()

	// 配置文件初始化
	// cfg 变量值从命令行 flag 传入，可以传值，比如 ./apiserver -c config.yaml，也可以为空，如果为空会默认读取 conf/config.yaml
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// Set gin mode
	/* 默认debug模式 */
	gin.SetMode(viper.GetString("runMode"))

	g := gin.New()

	// 加载路由&设置中间件
	router.Load(
		// cores
		g,

		// set middlwares
	)

	log.Infof("Listening and serving HTTP on %s\n", viper.GetString("addr"))
	// http.ListenAndServe()返回类型是error接口，调用Error()可以输出错误信息
	// g.Run()内部会调用http.ListenAndServe()
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// TODO 理论上wifi信号是一直在变化的，需要定时刷新参考点的数据，但是这有个问题就是参考点是人工添加的，能否实现自动更新，是否需要添加一些条件来实现自动更新
