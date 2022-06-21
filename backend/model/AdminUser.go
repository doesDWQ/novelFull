package model

import "gorm.io/gorm"

// gorm.Model 的定义
type AdminUser struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Passwrod string
	Email    string
	Token    string
}
