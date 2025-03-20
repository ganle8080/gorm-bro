package main

import (
	"fmt"
	"reflect"
)

func PrintNonZeroFields(obj interface{}) {
	// 获取对象的反射值
	v := reflect.ValueOf(obj)

	// 如果传递的是指针，则获取其指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 确保传递的是一个结构体
	if v.Kind() != reflect.Struct {
		fmt.Println("Expected a struct or a pointer to a struct")
		return
	}

	// 获取结构体的类型信息
	t := v.Type()

	// 遍历结构体的字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		// 如果字段的值不是零值，则打印字段名和值
		if !field.IsZero() {
			fmt.Printf("%s: %v\n", fieldName, field.Interface())
		}
	}

}

type MyStruct struct {
	Name   string
	Age    int
	Active bool
}

func main() {
	obj := MyStruct{
		Name:   "John",
		Age:    0,
		Active: true,
	}

	PrintNonZeroFields(obj)
}
