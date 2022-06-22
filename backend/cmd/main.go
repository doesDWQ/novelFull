package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/router"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements echo.Validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

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
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world!\n")
	})

	// 后台路由组
	router.AdminRoutes(e.Group("/admin"))

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port)))
}
