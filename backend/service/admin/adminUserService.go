package admin

import (
	"github.com.doesDWQ.novelFull/service/commonService"
)

type adminUserService struct {
	commonService.CommonService
}

// // 获取 adminUserService
// func NewAdminUserService() *adminUserService {
// 	return &adminUserService{
// 		CommonService: commonService.CommonService{
// 			// 获取原始model对象，用来表示需要查询的表
// 			Model: func() interface{} {
// 				return &model.AdminUser{}
// 			},

// 			// 列表时返回数据的model
// 			ListModel: func() interface{} {
// 				return &[]*model.AdminUser{}
// 			},

// 			// 详情和列表 返回数据的model
// 			SearchApiModel: func() interface{} {
// 				return &model.UserApi{}
// 			},

// 			// 校验添加数据的model
// 			AddRequestModel: func() interface{} {
// 				return &struct {
// 					UserName string `json:"username" validate:"required"`
// 					Passwrod string `json:"password" validate:"required"`
// 					Email    string `json:"email" validate:"required,email"`
// 				}{}
// 			},

// 			// 校验编辑数据的model
// 			EditRequestModel: func() interface{} {
// 				return &struct {
// 					UserName string `json:"username"`
// 					Passwrod string `json:"password"`
// 					Email    string `json:"email"`
// 				}{}
// 			},
// 		},
// 	}
// }
