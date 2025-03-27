package handler

import "gorm.io/gorm"

// GetHandlerFactory 根据 handlerName 返回对应的处理器工厂
func GetHandlerFactory(name string) (func(*gorm.DB) (interface{}, error), bool) {
	if name == "DemoTestHandler" {
		return func(db *gorm.DB) (interface{}, error) {
			return &DemoTestHandler{db: *db}, nil
		}, true
	}
	return nil, false
}
