package orm

import (
	"fmt"
	"gorm-bro/orm/handler"
	"reflect"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestHandler(t *testing.T) {
	// 初始化数据库连接
	var dsn string = "root:123456@tcp(1.95.212.179:3306)/demo_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("db: %v\n", db)

	// 动态执行方法
	handlerName := "DemoTestHandler"
	methodName := "AppendStr"

	// 获取处理器工厂
	factory, ok := handler.GetHandlerFactory(handlerName)
	if !ok {
		fmt.Printf("Handler not found: %s\n", handlerName)
		return
	}

	// 创建处理器实例
	instance, err := factory()
	if err != nil {
		fmt.Printf("Failed to create handler instance: %v\n", err)
		return
	}

	// 使用反射查找方法
	method := reflect.ValueOf(instance).MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("Method not found: %s\n", methodName)
		return
	}

	// 调用方法
	results := method.Call(nil) // 如果方法有参数，可以在这里传入参数切片
	if len(results) > 0 {
		fmt.Printf("Method result: %v\n", results[0].Interface())
	}

}

func TestSearch(t *testing.T) {
	map1 := map[string]interface{}{
		"table_name":  "demo_test",
		"search_type": "search",
		"order_by":    []string{},
		"columns":     []string{},
		"conditions": []map[string]interface{}{
			{
				"name":  "money",
				"type":  "eq",
				"value": 1.3,
			},
			{
				"name":  "age",
				"type":  "eq",
				"value": 3,
			},
		},
		"page": 1,
		"size": 20,
	}

	// 初始化数据库连接
	var dsn string = "root:123456@tcp(1.95.212.179:3306)/demo_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(Search(db, &map1))
}
