package dao

import (
	"fmt"

	"api.backend.xjco2913/util/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var (
	DB *gorm.DB // default db connection
) 

func init() {
	var err error

	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		config.Get("database.mysql.user"),
		config.Get("database.mysql.password"),
		config.Get("database.mysql.host"),
		config.Get("database.mysql.port"),
		config.Get("database.mysql.databaseName"),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}