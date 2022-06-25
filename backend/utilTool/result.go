package utilTool

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// 成功的返回结构
func Success(c echo.Context, data map[string]interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"data":   data,
		"status": "1",
		"msg":    "",
	})
}

// 失败的返回结构
func Error(c echo.Context, err error) error {
	e := c.JSON(http.StatusOK, echo.Map{
		"data":   nil,
		"status": "0",
		"msg":    err.Error(),
	})
	if e != nil {
		return e
	}
	return err
}
