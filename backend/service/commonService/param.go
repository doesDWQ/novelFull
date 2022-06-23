package commonService

import (
	"fmt"
	"strconv"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/validate"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// context 参数处理定义

// Paginate 分页
func (common *CommonService) Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetTokenInfo 获取token信息
func (common *CommonService) GetTokenInfo(c echo.Context) (uint, *jwt.Token) {
	i := c.Get(config.Config.Jwt.ContextKey)
	token := i.(*jwt.Token)

	claims := token.Claims.(*AdminJwtCustomClaims)

	return claims.UserId, token
}

// DealParam 根据输入规则map 处理参数，返回处理校验后的map型数据
func (common *CommonService) DealParam(c echo.Context, rules map[string]interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if err := c.Bind(params); err != nil {
		return nil, fmt.Errorf("bind err: %s", err)
	}

	fmt.Printf("%#v\n", params)

	if err := c.Validate(&validate.MapValidate{
		Value: params,
		Rules: common.AddRules,
	}); err != nil {
		return nil, common.Error(c, fmt.Sprintf("validate err: %s", err))
	}

	// 根据addRules取出所有需要保存的字段
	data := make(map[string]interface{})
	for k, _ := range rules {
		if value, exists := params[k]; exists {
			data[k] = value
		}
	}

	return data, nil
}

// 获取分页信息
func (common *CommonService) getPageInfo(c echo.Context) (*pageInfo, error) {
	pageString := c.FormValue("page")
	pageSizeString := c.FormValue("pageSize")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		return nil, err
	}
	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		return nil, err
	}

	return &pageInfo{
		Page:     page,
		PageSize: pageSize,
	}, nil
}
