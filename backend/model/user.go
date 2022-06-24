package model

import "gorm.io/gorm"

// gorm.Model 的定义
type User struct {
	gorm.Model
	UserName string
	Passwrod string
	Email    string
	Token    string
}

// 查询模型定义
type UserApi struct {
	gorm.Model
	UserName string
	Email    string
}
