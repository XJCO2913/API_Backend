package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var (
	DB *gorm.DB // default db connection
) 

func init() {
	var err error
	dsn := "xiaofei:2021110003@tcp(127.0.0.1:3306)/API_XJCO2913?parseTime=true"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}