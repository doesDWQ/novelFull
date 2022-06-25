package commonService

import (
	"errors"

	"github.com.doesDWQ.novelFull/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RegisterRouterFunc func(g *echo.Group) interface{}

// 结构体定义
type (
	Route struct {
		// 请求方法
		RequestFunc func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		// 请求子路径
		Path string
		// 实际请求方法
		Func func(c echo.Context) error
	}

	// 内部service表示service初始化是必须要考虑的参数
	InnerService struct {
		Db *gorm.DB
		// 获取原始model的方法
		Model func() model.CommonModelInterface
		// 添加时的请求参数model
		AddRules map[string]interface{}
		// 编辑时的请求参数model
		EditRules map[string]interface{}
		// 注册路由
		Routes []*Route
		// 当前所属路由组
		Group *echo.Group
		// 是否自动注册默认的增删改查等接口
		AutoRegister bool
	}

	// 这个结构体中的成员必须传递过来，并且必须大写开头，而且值类型必须是指针(指针方便判空)
	Service struct {
		innerService *InnerService
	}

	ServiceInterface interface {
		GetRoutes() ([]*Route, error)
	}

	JwtCustomClaims struct {
		UserId   uint   `json:user_id`
		UserName string `json:"user_name"`
		Admin    bool   `json:"admin"`
		jwt.StandardClaims
	}
)

// 初始化内部service
func (common *Service) InitInnerService(innerService *InnerService) {
	common.innerService = innerService
}

// 获取默认注册的路由
func (common *Service) getDefaultRoutes() []*Route {
	return []*Route{
		// 添加
		{
			RequestFunc: common.innerService.Group.POST,
			Path:        "",
			Func:        common.AddAPi,
		},
		// 编辑
		{
			RequestFunc: common.innerService.Group.PUT,
			Path:        "/:id",
			Func:        common.EditApi,
		},
		// 详情
		{
			RequestFunc: common.innerService.Group.GET,
			Path:        "/:id",
			Func:        common.DetailApi,
		},
		// 删除
		{
			RequestFunc: common.innerService.Group.DELETE,
			Path:        "/:id",
			Func:        common.DeleteApi,
		},
		// 列表
		{
			RequestFunc: common.innerService.Group.GET,
			Path:        "",
			Func:        common.ListApi,
		},
	}
}

// 获取路由数组，路由只在启动的时候注册一次
func (common Service) GetRoutes() ([]*Route, error) {
	// 判断是否初始化过内部service
	if common.innerService == nil {
		return nil, errors.New("未初始化内部service")
	}

	if common.innerService.AutoRegister {
		// 是否注册默认路由
		return append(common.getDefaultRoutes(), common.innerService.Routes...), nil
	} else {
		return common.innerService.Routes, nil
	}
}
