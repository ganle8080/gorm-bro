package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	var dsn string = "root:123456@tcp(1.95.212.179:3306)/demo_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	data := []map[string]interface{}{}
	db.Raw("desc demo_test").Scan(&data)
	fmt.Printf("data: %v\n", data[0])
}

/*

SELECT
    COLUMN_NAME,
    COLUMN_COMMENT
FROM
    information_schema.COLUMNS
WHERE
    TABLE_SCHEMA = 'demo_go'
    AND TABLE_NAME = 'demo_test';
*/
