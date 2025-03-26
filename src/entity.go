package src

import (
	"gorm.io/gorm"
)

type DemoTest struct {
	gorm.Model
	Name  string  `gorm:"column:name;type:varchar(100);not null"`       // 名称字段，最大长度100，不能为空
	Age   int     `gorm:"column:age;type:int;default:18"`               // 年龄字段，默认值为18
	Money float64 `gorm:"column:money;type:decimal(10,2);default:0.00"` // 金额字段，精度为10，小数位为2，默认值为0.00
}

func (d *DemoTest) TableName() string {
	return "demo_test"
}
