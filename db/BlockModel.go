package db

import (
	"github.com/jinzhu/gorm"
)

type Blok struct {
	gorm.Model
	PreHash       string  `gorm:"column:prehash"`
	Data          string  `gorm:"column:data"`
	Timestamp     int  `gorm:"column:timestamp"`
	TargetBits    int `gorm:"column:targetBits"`
	Noce         int  `gorm:"column:noce"`

}

