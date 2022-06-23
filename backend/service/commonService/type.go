package commonService

import (
	"github.com.doesDWQ.novelFull/types"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type RegisterRouterFunc func(g *echo.Group) interface{}

// 结构体定义
type (
	CommonService struct {
		// 获取原始model的方法
		Model func() interface{}
		// 列表model
		ListModel func() interface{}
		// 查询model
		SearchApiModel func() interface{}
		// 添加时的请求参数model
		AddRules map[string]interface{}
		// 编辑时的请求参数model
		EditRules map[string]interface{}
		// 注册路由
		Routes []*types.Route
		// 是否自动注册默认的增删改查等接口
		AutoRegister bool
	}

	AdminJwtCustomClaims struct {
		UserId   uint   `json:user_id`
		UserName string `json:"user_name"`
		Admin    bool   `json:"admin"`
		jwt.StandardClaims
	}

	pageInfo struct {
		Page     int
		PageSize int
	}
)
