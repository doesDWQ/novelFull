package commonService

import (
	"errors"
	"fmt"
	"strconv"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/types"
	"github.com.doesDWQ.novelFull/validate"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// context 参数处理定义

// GetTokenInfo 获取token信息
func (common *Service) GetTokenInfo(c echo.Context) (uint, *jwt.Token) {
	i := c.Get(config.Config.Jwt.ContextKey)
	token := i.(*jwt.Token)

	claims := token.Claims.(*JwtCustomClaims)

	return claims.UserId, token
}

// DealParam 根据输入规则map 处理参数，返回处理校验后的map型数据
func (common *Service) DealParam(c echo.Context, rules map[string]interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if err := c.Bind(&params); err != nil {
		return nil, fmt.Errorf("bind err: %s", err)
	}

	fmt.Printf("params: %#v,rules: %#v \n", params, rules)

	err := c.Validate(&validate.MapValidate{
		Value: params,
		Rules: rules,
	})
	fmt.Println("v2.error:", err)
	if err != nil {
		// 发送错误的相应
		return nil, common.Error(c, fmt.Errorf("validate err: %s", err))
	}

	fmt.Println("v3.error:", err)

	data := make(map[string]interface{})
	for k, _ := range rules {
		if value, exists := params[k]; exists {
			data[k] = value
		}
	}

	return data, nil
}

// 获取分页信息
func (common *Service) GetPageInfo(c echo.Context) (*types.PageInfo, error) {
	params := make(map[string]interface{})
	if err := c.Bind(&params); err != nil {
		return nil, fmt.Errorf("bind err: %s", err)
	}

	if params["page"] == nil || params["pageSize"] == nil {
		return nil, errors.New("not get params: page and pageSize")
	}

	page, err := strconv.Atoi(params["page"].(string))
	if err != nil {
		return nil, err
	}
	pageSize, err := strconv.Atoi(params["pageSize"].(string))
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v", params)

	return &types.PageInfo{
		Page:     page,
		PageSize: pageSize,
	}, nil
}
