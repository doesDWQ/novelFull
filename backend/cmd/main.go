package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/service/admin"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 初始化配置
	err := config.InitConfig("config/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	// 初始化数据库
	err = db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world!\n")
	})

	l := admin.LoginService{}
	// 后台登录接口
	e.POST("/login", l.Login)

	// 后台路由组
	AdminGroup := e.Group("/admin")
	AdminGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
		Claims:     &admin.AdminJwtCustomClaims{},
		SuccessHandler: func(c echo.Context) {
			// token校验成功后的操作，查询该token对该用户是否有效

		},
	}))

	adminUser := admin.AdminUserService{}
	// 新增
	AdminGroup.POST("/users", adminUser.Add)
	// 编辑
	AdminGroup.PUT("/users/:id", adminUser.Add)
	// 详情
	AdminGroup.GET("/users/:id", adminUser.Detail)
	// 列表
	AdminGroup.GET("/users", adminUser.List)
	// 删除
	AdminGroup.DELETE("/users/:id", adminUser.Delete)

	AdminGroup.Any("/hello", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world!\n")
	})

	// 退出接口，登录的时候需要将token保存到用户表，就能单点登录了，失效则直接删除用户的token即可
	AdminGroup.POST("/loginout", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world!\n")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port)))
}
