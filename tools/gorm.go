package tools

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// List 查询符合object对象的列表
func List(db *gorm.DB, object interface{}) {
	// 如果object是值对象，获取其指针
	objType := reflect.TypeOf(object)
	if objType.Kind() != reflect.Ptr {
		object = reflect.New(objType).Interface()
	}

	// 获取object的结构体名称作为表名
	tableName := reflect.TypeOf(object).Elem().Name()

	// 判断object是否实现了TableName方法
	if v, ok := object.(interface {
		TableName() string
	}); ok {
		tableName = v.TableName()
	}

	fmt.Printf("tableName: %v\n", toSnakeCase(tableName))
}

// 将驼峰命名转换成为蛇形命名
func toSnakeCase(name string) string {
	result := ""
	for i, char := range name {
		if char >= 'A' && char <= 'Z' {
			// 如果不是第一个字符，前面加下划线
			if i > 0 {
				result += "_"
			}
			// 转为小写
			result += string(char + ('a' - 'A'))
		} else {
			result += string(char)
		}
	}
	return result
}
