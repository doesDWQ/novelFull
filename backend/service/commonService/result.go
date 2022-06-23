package commonService

import (
	"github.com.doesDWQ.novelFull/utilTool"
	"github.com/labstack/echo/v4"
)

// 成功的返回结构
func (common *CommonService) Success(c echo.Context, data map[string]interface{}) error {
	return utilTool.Success(c, data)
}

// 失败的返回结构
func (common *CommonService) Error(c echo.Context, errMsg string) error {
	return utilTool.Error(c, errMsg)
}
