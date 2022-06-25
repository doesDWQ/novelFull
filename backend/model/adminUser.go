package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AdminUser struct {
	gorm.Model
	CommonModel

	UserName string `gorm:"unique"`
	Pwd      string
	Email    string
	Token    string
}

// 查询模型定义, 展示给前端的数据
type SearchAdminUserApi struct {
	gorm.Model
	UserName string
	Email    string
}

// 标准都应该使用这个获取到model
func NewAdminUser() *AdminUser {
	modelPtr := &AdminUser{}

	// 设置当前model
	modelPtr.SetCurrentModel(modelPtr)

	// 设置单个查询model
	modelPtr.SetModelFunc(func() interface{} {
		return &SearchAdminUserApi{}
	})

	// 设置列表查询model
	modelPtr.SetListModelFunc(func() interface{} {
		return &[]*SearchAdminUserApi{}
	})

	return modelPtr
}

// 根据用户名查找用户
func (au *AdminUser) FindByUsernamePwd(tx *gorm.DB, userName, pwd string) (searchModel interface{}, err error) {
	searchModel, err = au.GetModelFunc()
	if err != nil {
		return
	}
	err = tx.Model(au).Debug().Where(map[string]interface{}{
		"user_name": userName,
		"pwd":       pwd,
	}).Scan(searchModel).Error

	after := searchModel.(*SearchAdminUserApi)
	if after.ID == 0 {
		err = errors.New("record is not record")
		return
	}

	fmt.Printf("searchModel:%#v, \n", searchModel)

	return
}

// 根据用户名查找用户
func (au *AdminUser) FindByIdToken(tx *gorm.DB, id, token interface{}) (searchModel interface{}, err error) {
	searchModel, err = au.GetModelFunc()
	if err != nil {
		return
	}
	err = tx.Model(au).Debug().Where(map[string]interface{}{
		"id":    id,
		"token": token,
	}).Scan(searchModel).Error

	after := searchModel.(*SearchAdminUserApi)
	if after.ID == 0 {
		err = errors.New("record is not record")
		return
	}

	fmt.Printf("searchModel:%#v, \n", searchModel)

	return
}

// 更新token
func (au *AdminUser) UpdateTokenById(tx *gorm.DB, id uint, token string) error {
	return au.UpdateById(tx, int(id), map[string]interface{}{
		"token": token,
	})
}
