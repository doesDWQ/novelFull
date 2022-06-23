package router

import (
	"errors"
	"fmt"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service/admin"
	"github.com.doesDWQ.novelFull/service/commonService"
	"github.com.doesDWQ.novelFull/types"
	"github.com.doesDWQ.novelFull/utilTool"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var controllers = map[string]map[string]func(e *echo.Group) types.Service{
	"/admin": {
		"/basic": admin.NewLoginService,
	},
}

// 跳过权限校验的接口
var skipApis = map[string]struct{}{
	"/admin/basic/login": {},
}

// 注册路由
func RegisterRoutes(e *echo.Echo) {

	// 后端路由
	g := e.Group("/",
		// 注册两个权限校验中间件
		middleware.JWTWithConfig(middleware.JWTConfig{
			// 跳过token校验接口
			Skipper:    skipVerify,
			Claims:     &commonService.AdminJwtCustomClaims{},
			SigningKey: []byte(config.Config.Jwt.Secret),
			ContextKey: config.Config.Jwt.ContextKey,
		}),
		verifyToken,
	)

	for firstPath, secondRouter := range controllers {
		// 注册第一级路由组
		nextGroup := g.Group(firstPath)
		// 注册第二层路由组
		for secondPath, sonRouterFunc := range secondRouter {
			sonGroup := nextGroup.Group(secondPath)
			service := sonRouterFunc(sonGroup)

			// 注册第三级路由
			for _, lastRoute := range service.GetRoutes() {
				// 这里还可以返回下一级别的router，但是不处理了
				lastRoute.RequestFunc(lastRoute.Path, lastRoute.Func)
				// 处理是否跳过
				if lastRoute.SkipVerify {
					path := fmt.Sprintf("%s%s%s", firstPath, secondPath, lastRoute.Path)
					skipApis[path] = struct{}{}
				}
			}

			// 处理是否自动注册，增删改查列表编辑接口
			if service.GetAutoRegister() {
				sonGroup.POST("/", service.AddAPi)
				sonGroup.PUT("/:id", service.EditApi)
				sonGroup.GET("/:id", service.DetailApi)
				sonGroup.GET("/", service.ListApi)
				sonGroup.DELETE("/:id", service.DeleteApi)
			}
		}
	}
}

// 校验token是否存在
func verifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if skipVerify(c) {
			// 过权限校验
			return next(c)
		}

		// 校验token
		cs := commonService.CommonService{}
		userId, token := cs.GetTokenInfo(c)
		user := &model.AdminUser{
			Model: gorm.Model{
				ID: userId,
			},
		}

		err := db.Db.Where(user).First(user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utilTool.Error(c, "用户id不存在")
			}
			return err
		}

		// token为空表示未登录，token和表里面的token没有对照上也是登录错误
		if user.Token == "" && user.Token != token.Raw {
			// token校验不存在则报错
			return utilTool.Error(c, "用户信息校验失败")
		}

		return next(c)
	}
}

// 跳过权限校验的路由
func skipVerify(c echo.Context) bool {
	url := c.Request().URL.String()

	if _, exists := skipApis[url]; exists {
		return true
	}

	return false
}
