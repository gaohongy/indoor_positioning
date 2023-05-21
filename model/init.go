package model

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/zxmrlc/log"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	Mysql *gorm.DB
}

var DB *Database

func (db *Database) Init() {
	DB = &Database{
		Mysql: GetMysqlDB(),
	}
}

func (db *Database) Close() {
	DB.Mysql.Close()
}

func GetMysqlDB() *gorm.DB {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s", //
		username,
		password,
		addr,
		name,
		// viper.GetString("db.encode"),
		true,
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
	} else {
		log.Infof("Database connection success. Database name: %s", name)
	}

	// 设置数据库属性
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.SingularTable(true) // 在gorm中，默认的表名为定义的结构体名首字母大写转小写，然后变为复试，该选项为true则禁止变为复数
	db.LogMode(viper.GetBool("gormlog"))
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}
