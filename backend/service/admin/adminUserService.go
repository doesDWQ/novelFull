package admin

import (
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service/commonService"
	"github.com/labstack/echo/v4"
)

type adminUserService struct {
	commonService.Service
}

// 获取 adminUserService
func NewAdminUserService(g *echo.Group) commonService.ServiceInterface {
	service := adminUserService{}
	service.InitInnerService(&commonService.InnerService{
		Db: db.Db,
		Model: func() model.CommonModelInterface {
			return model.NewAdminUser()
		},
		AddRules: map[string]interface{}{
			"user_name": "required,min=3,max=14",
			"pwd":       "required,min=6,max=22",
			"email":     "required,email",
		},
		EditRules: map[string]interface{}{
			"user_name": "min=3,max=14",
			"pwd":       "min=6,max=22",
			"email":     "email",
		},
		Routes:       []*commonService.Route{},
		Group:        g,
		AutoRegister: true,
	})
	return &service
}
