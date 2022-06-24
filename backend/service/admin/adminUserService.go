package admin

import (
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service/commonService"
	"github.com.doesDWQ.novelFull/types"
	"github.com/labstack/echo/v4"
)

type adminUserService struct {
	commonService.Service
}

// 获取 adminUserService
func NewAdminUserService(g *echo.Group) types.Service {
	return &adminUserService{
		Service: commonService.Service{
			Model: func() interface{} {
				return model.AdminUser{}
			},
			ListModel: func() interface{} {
				return []*model.AdminUserApi{}
			},
			SearchApiModel: func() interface{} {
				return model.AdminUserApi{}
			},
			AddRules: map[string]interface{}{
				"user_name": "required,min=3,max=14",
				"pwd":       "required,min=6,max=22",
				"email":     "required,email",
			},
			EditRules: map[string]interface{}{
				"user_name": "min=3,max=14",
				"pwd":       "min=6,max=22",
				"email":     "required,email",
			},
			Routes:       []*types.Route{},
			Group:        g,
			AutoRegister: true,
		},
	}
}
