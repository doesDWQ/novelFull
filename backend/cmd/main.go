package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/router"
	"github.com.doesDWQ.novelFull/validate"

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
	// 校验器
	e.Validator = validate.NewCustomValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world!\n")
	})

	// 注册路由
	err = router.RegisterRoutes(e)
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port)))
}
