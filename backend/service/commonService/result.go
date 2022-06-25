package commonService

import (
	"github.com.doesDWQ.novelFull/utilTool"
	"github.com/labstack/echo/v4"
)

// 成功的返回结构
func (common *Service) Success(c echo.Context, data map[string]interface{}) error {
	return utilTool.Success(c, data)
}

// 失败的返回结构
func (common *Service) Error(c echo.Context, err error) error {
	if err == nil {
		return nil
	}
	return utilTool.Error(c, err)
}
