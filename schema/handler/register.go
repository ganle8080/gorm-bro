package handler

import (
	"fmt"
	"sync"
)

// HandlerFactory 是一个工厂函数类型，用于创建处理程序实例
type HandlerFactory func(args ...interface{}) (interface{}, error)

var (
	factoryRegistry = make(map[string]HandlerFactory)
	registryLock    sync.RWMutex
)

// RegisterHandler 注册一个新的处理程序工厂
func RegisterHandler(name string, factory HandlerFactory) {
	registryLock.Lock()
	defer registryLock.Unlock()
	factoryRegistry[name] = factory

	fmt.Printf("factoryRegistry: %v\n", factoryRegistry)
}

// GetHandlerFactory 根据名称获取处理程序工厂
func GetHandlerFactory(name string) (HandlerFactory, bool) {
	registryLock.RLock()
	defer registryLock.RUnlock()
	factory, exists := factoryRegistry[name]
	return factory, exists
}
