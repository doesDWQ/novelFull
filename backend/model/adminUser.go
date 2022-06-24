package model

import "gorm.io/gorm"

// gorm.Model 的定义
type AdminUser struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Pwd      string
	Email    string
	Token    string
}

// 查询模型定义, 展示给前端的数据
type AdminUserApi struct {
	gorm.Model
	UserName string
	Email    string
}
