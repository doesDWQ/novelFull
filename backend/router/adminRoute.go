package router

import (
	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/service"
	"github.com.doesDWQ.novelFull/service/admin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AdminRoutes(g *echo.Group) error {

	loginService := admin.NewLoginService()
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		// 跳过token校验接口
		Skipper:    loginService.SkipVerify,
		Claims:     &service.AdminJwtCustomClaims{},
		SigningKey: []byte(config.Config.Jwt.Secret),
		ContextKey: config.Config.Jwt.ContextKey,
	}))
	// 校验token
	g.Use(loginService.VerifyToken)

	// 登录接口
	g.POST("/login", loginService.Login)
	// 退出
	g.POST("/loginout", loginService.LoginOut)
	// 获取用户信息
	g.GET("/userInfo", loginService.AdminUserInfo)

	// 后端用户管理
	adminUserService := admin.NewAdminUserService()
	g.POST("/users", adminUserService.Add)
	g.PUT("/users/:id", adminUserService.Edit)
	g.GET("/users/:id", adminUserService.Detail)
	g.GET("/users", adminUserService.List)
	g.DELETE("/users/:id", adminUserService.Delete)

	return nil
}
