package src

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DemoTest struct {
	gorm.Model
	Name  string
	Age   int
	Money float64
	Desc  datatypes.JSON
}

func (d *DemoTest) TableName() string {
	return "demo_test_giao"
}
