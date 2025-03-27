package handler

import (
	"fmt"

	"gorm.io/gorm"
)

// DemoTestHandler 示例处理器
type DemoTestHandler struct {
	db gorm.DB
}

// AppendStr 示例方法
func (h *DemoTestHandler) AppendStr(money interface{}) string {
	fmt.Printf("h.db: %v\n", h.db)
	return "Hello, World!" + fmt.Sprintf("%v", money)
}

// GetHandlerFactory 根据 handlerName 返回对应的处理器工厂
func GetHandlerFactory(name string) (func(*gorm.DB) (interface{}, error), bool) {
	if name == "DemoTestHandler" {
		return func(db *gorm.DB) (interface{}, error) {
			return &DemoTestHandler{db: *db}, nil
		}, true
	}
	return nil, false
}
