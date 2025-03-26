package main

import (
	"fmt"
	"gorm-bro/schema/handler"
	"reflect"
	"strconv"

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

	// map1 := map[string]interface{}{
	// 	"table_name":  "demo_test",
	// 	"search_type": "list",
	// 	"order_by":    []string{},
	// 	"columns":     []string{},
	// 	"conditions": []map[string]interface{}{
	// 		{
	// 			"name":  "name",
	// 			"type":  "eq",
	// 			"value": "helloworld",
	// 		},
	// 	},
	// 	"page": 1,
	// 	"size": 20,
	// }

	// tools.Search(map1, db)

	handlerName := "DemoTestHandler"
	methodName := "AppendStr"

	factory, ok := handler.GetHandlerFactory(handlerName)
	if !ok {
		fmt.Println("factory not found")
		return
	}

	instance, err := factory(db)

	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return
	}

	// 使用反射获取实例的类型
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	method := instanceValue.MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("Method '%s' not found on handler '%s'\n", methodName, handlerName)
		return
	}

	var methodArgs []reflect.Value
	// 如果方法需要参数，这里需要根据方法签名进行转换
	// 例如，假设方法有一个 string 参数
	methodArgs = append(methodArgs, reflect.Value{})
	str := method.Call(methodArgs)
	fmt.Printf("str: %v\n", str)
}

func convertParam(paramStr string, targetType reflect.Type) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(paramStr), nil
	case reflect.Int:
		i, err := strconv.Atoi(paramStr)
		if err != nil {
			return reflect.Zero(targetType), fmt.Errorf("invalid integer value: %s", paramStr)
		}
		return reflect.ValueOf(i), nil
	// 添加更多类型转换
	default:
		return reflect.Zero(targetType), fmt.Errorf("unsupported parameter type: %s", targetType)
	}
}
