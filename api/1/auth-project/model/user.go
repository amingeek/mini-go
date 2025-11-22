package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"unique;type:varchar(255);not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
}
