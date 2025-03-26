package handler

import (
	"fmt"

	"gorm.io/gorm"
)

type DemoTestHandler struct {
	db *gorm.DB
}

func NewDemoTestHandler(db *gorm.DB) *DemoTestHandler {
	return &DemoTestHandler{db: db}
}

func (h *DemoTestHandler) AppendStr(str interface{}) string {
	return "$ " + fmt.Sprintf("%s", str)
}

// demoTestHandlerFactory 是 DemoTestHandler 的工厂函数
func demoTestHandlerFactory(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("DemoTestHandler requires at least one argument: Name")
	}
	db, ok := args[0].(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid type for Name, expected string")
	}
	return &DemoTestHandler{db: db}, nil
}

func init() {
	RegisterHandler("DemoTestHandler", demoTestHandlerFactory)
}
