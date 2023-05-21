package config

import (
	"indoor_positioning/model"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zxmrlc/log"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包
	// 日志初始化需要读取配置文件，故此初始化位于配置文件初始化之后
	c.initLog()

	// 初始化数据库
	model.DB.Init()
	// https://blog.csdn.net/HYZX_9987/article/details/116698754 这里查询到是数据库提前关闭了
	// defer model.DB.Close() // 延迟关闭数据库，defer确保资源能够被关闭

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" { // 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else { // 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("config")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML

	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

// 监控配置文件变化并热加载程序
// 热更新是指：可以不重启 API 进程，使 API 加载最新配置项的值
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
