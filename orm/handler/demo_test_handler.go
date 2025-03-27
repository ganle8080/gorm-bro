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
