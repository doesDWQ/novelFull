package commonService

import (
	"errors"

	"github.com.doesDWQ.novelFull/types"
)

// 获取路由数组
func (common *Service) GetRoutes() ([]*types.Route, error) {
	// 检查默认参数是否已经设置
	if common.Model == nil ||
		common.ListModel == nil ||
		common.SearchApiModel == nil ||
		common.AddRules == nil ||
		common.EditRules == nil ||
		common.Routes == nil ||
		common.Group == nil {

		return nil, errors.New("service 的必须参数：Model,ListModel, SearchApiModel, AddRules, EditRules, Routes, Group 未初始化")
	}

	if common.registedRoutes {
		return common.allRoutes, nil
	}

	if common.defaultRoutes == nil {
		common.defaultRoutes = []*types.Route{
			{
				RequestFunc: common.Group.POST,
				Path:        "",
				Func:        common.AddAPi,
				SkipVerify:  false,
			},
			{
				RequestFunc: common.Group.PUT,
				Path:        "/:id",
				Func:        common.EditApi,
				SkipVerify:  false,
			},
			{
				RequestFunc: common.Group.GET,
				Path:        "/:id",
				Func:        common.DetailApi,
				SkipVerify:  false,
			},
			{
				RequestFunc: common.Group.DELETE,
				Path:        "/:id",
				Func:        common.DeleteApi,
				SkipVerify:  false,
			},
			{
				RequestFunc: common.Group.GET,
				Path:        "",
				Func:        common.ListApi,
				SkipVerify:  false,
			},
		}
	}

	// 注册默认路由
	if common.AutoRegister {
		common.allRoutes = append(common.allRoutes, common.defaultRoutes...)
	}

	// 注册声明路由
	common.allRoutes = append(common.allRoutes, common.Routes...)

	common.registedRoutes = true

	return common.allRoutes, nil
}
