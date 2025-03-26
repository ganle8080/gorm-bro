package main

import (
	"fmt"
	"gorm-bro/tools"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var dsn string = "root:123456@tcp(127.0.0.1:3306)/demo_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("db: %v\n", db)

	map1 := map[string]interface{}{
		"table_name":  "demo_test",
		"search_type": "list",
		"order_by":    []string{},
		"columns":     []string{},
		"conditions": []map[string]interface{}{
			{
				"name":  "name",
				"type":  "eq",
				"value": "helloworld",
			},
		},
		"page": 1,
		"size": 20,
	}

	tools.Search(map1, db)
}
