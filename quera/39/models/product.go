package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string `gorm:"type:varchar(255);not null"`
	Count int64  `gorm:"type:int"`
	Price int64  `gorm:"type:bigint"`
}
