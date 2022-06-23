package types

import (
	"github.com/labstack/echo/v4"
)

type Route struct {
	// 请求方法
	RequestFunc func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	// 请求子路径
	Path string
	// 实际请求方法
	Func func(c echo.Context) error
	// 该接口是否跳过权限校验
	SkipVerify bool
}

type Service interface {
	GetRoutes() []*Route
	GetAutoRegister() bool

	AddAPi(c echo.Context) error
	DetailApi(c echo.Context) error
	DeleteApi(c echo.Context) error
	ListApi(c echo.Context) error
	EditApi(c echo.Context) error
}
