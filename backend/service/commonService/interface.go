package commonService

import "github.com.doesDWQ.novelFull/types"

// 获取路由数组
func (common *CommonService) GetRoutes() []*types.Route {
	return common.Routes
}

// 获取是否注册属性值
func (common *CommonService) GetAutoRegister() bool {
	return common.AutoRegister
}
