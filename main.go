package main

import (
	"gorm-bro/src"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var dsn string = "root:123456@tcp(1.95.212.179:3306)/demo_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	demoTest := src.DemoTest{}

	db.AutoMigrate(&demoTest)

}
