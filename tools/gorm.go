package tools

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// List 查询符合object对象的列表
func List(db *gorm.DB, object interface{}, a ...string) {
	// 如果object是值对象，获取其指针
	objType := reflect.TypeOf(object)
	if objType.Kind() != reflect.Ptr {
		object = reflect.New(objType).Interface()
	}

	// 如果a有值，打印第一个值
	if len(a) > 0 {
		fmt.Printf("a[0]: %v\n", a)
	}

	fmt.Printf("tableName: %v\n", getTableName(object))

	// 获取object的字段名，将其转换为蛇形命名
	fmt.Printf("columns: %v\n", getColumns(object))
}

// select demo_test.id ... from demo_test where

func getColumns(object interface{}) []string {
	// 获取对象的类型
	objType := reflect.TypeOf(object)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem() // 如果是指针，获取其指向的元素类型
	}

	// 如果不是结构体，返回空切片
	if objType.Kind() != reflect.Struct {
		return []string{}
	}

	columns := []string{}
	// 遍历结构体的字段
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		// 检查是否有gorm:"-"标签
		if tag, ok := field.Tag.Lookup("gorm"); ok && tag == "-" {
			continue // 忽略带有gorm:"-"标签的字段
		}

		// 如果字段是嵌套的匿名结构体，递归处理
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedColumns := getColumns(reflect.New(field.Type).Interface())
			columns = append(columns, embeddedColumns...)
			continue
		}

		// 获取字段名并转换为蛇形命名
		columnName := toSnakeCase(field.Name)
		columns = append(columns, columnName)
	}

	return columns
}
func getTableName(object interface{}) string {
	// 获取object的结构体名称作为表名
	tableName := reflect.TypeOf(object).Elem().Name()
	// 判断object是否实现了TableName方法
	if v, ok := object.(interface {
		TableName() string
	}); ok {
		tableName = v.TableName()
	}
	return toSnakeCase(tableName)
}

// 将驼峰命名转换成为蛇形命名
// 将驼峰命名转换成为蛇形命名
func toSnakeCase(name string) string {
	// 特殊处理：如果字段名是 "ID"，直接返回 "id"
	if name == "ID" {
		return "id"
	}

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
