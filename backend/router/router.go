package router

import (
	"fmt"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service/admin"
	"github.com.doesDWQ.novelFull/service/commonService"
	"github.com.doesDWQ.novelFull/service/loginService"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GetServiceFunc func(e *echo.Group) commonService.ServiceInterface
type FirstRoute struct {
	secondRoutes map[string]GetServiceFunc
	middlewares  []echo.MiddlewareFunc
}

// 跳过权限校验的接口
var skipApis = map[string]struct{}{
	"/admin/basic/adminLogin": {},
}

// 注册路由
func RegisterRoutes(e *echo.Echo) error {

	controllers := map[string]FirstRoute{
		"/admin": {
			secondRoutes: map[string]GetServiceFunc{
				"/basic":     loginService.NewLoginService,
				"/adminUser": admin.NewAdminUserService,
			},
			middlewares: []echo.MiddlewareFunc{
				// 注册两个权限校验中间件
				middleware.JWTWithConfig(middleware.JWTConfig{
					// 跳过token校验接口
					Skipper:    skipVerify,
					Claims:     &commonService.JwtCustomClaims{},
					SigningKey: []byte(config.Config.Jwt.Secret),
					ContextKey: config.Config.Jwt.ContextKey,
				}),
				verifyToken,
			},
		},
	}

	for firstPath, firstRoute := range controllers {
		// 注册第一级路由组
		nextGroup := e.Group(firstPath, firstRoute.middlewares...)
		// fmt.Printf("一级路径路径：%s\n", firstPath)
		// 注册第二层路由组
		for servicePath, getServiceFunc := range firstRoute.secondRoutes {
			// fmt.Printf("二级路径路径：%s\n", secondPath)
			// 注册组路由
			serviceGroup := nextGroup.Group(servicePath)
			service := getServiceFunc(serviceGroup)

			// 注册第三级路由
			thirdRoutes, err := service.GetRoutes()
			if err != nil {
				return fmt.Errorf(
					"path:%s: %s",
					fmt.Sprintf("%s%s", firstPath, servicePath), err.Error(),
				)
			}
			for _, lastRoute := range thirdRoutes {
				route := lastRoute.RequestFunc(lastRoute.Path, lastRoute.Func)
				path := fmt.Sprintf("%s%s%s", firstPath, servicePath, lastRoute.Path)
				fmt.Printf("当前注册路由: %s:%s\n", route.Method, path)
			}
		}
	}

	return nil
}

// 校验token是否存在
func verifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if skipVerify(c) {
			// 过权限校验
			return next(c)
		}

		// 校验token
		cs := commonService.Service{}
		userId, token := cs.GetTokenInfo(c)
		au := model.NewAdminUser()
		_, err := au.FindByIdToken(db.Db, userId, token.Raw)
		if err != nil {
			return err
		}

		return next(c)
	}
}

// 跳过权限校验的路由
func skipVerify(c echo.Context) bool {
	url := c.Request().URL.String()

	// fmt.Printf("当前调用到了，url:%s\n", url)
	// fmt.Printf("skipApis:%#v\n", skipApis)
	if _, exists := skipApis[url]; exists {
		// fmt.Println("当前调用到了，需要跳过")
		return true
	}

	return false
}
