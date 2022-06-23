package commonService

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Add 新增
func (common *CommonService) AddAPi(c echo.Context) error {
	if data, err := common.DealParam(c, common.AddRules); err != nil {
		return err
	} else {
		return common.add(data)
	}
}

// Detail 详情
func (common *CommonService) DetailApi(c echo.Context) error {

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("detail strconv error: %s", err)
	}

	data, err := common.detail(idInt)
	if err != nil {
		return err
	}
	err = common.Success(c, map[string]interface{}{
		"user": data,
	})

	if err != nil {
		return err
	}

	return nil
}

// Delete 删除
func (common *CommonService) DeleteApi(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	err = common.deleteById(idInt)
	if err != nil {
		return err
	}

	return nil
}

// List 列表
func (common *CommonService) ListApi(c echo.Context) error {
	// 获取分页信息
	pageInfo, err := common.getPageInfo(c)
	if err != nil {
		return err
	}

	listData, count, query := common.getListQuery(pageInfo)
	err = query.Error
	if err != nil {
		return err
	}

	err = common.Success(c, map[string]interface{}{
		"users": listData,
		"cnt":   count,
	})

	if err != nil {
		return err
	}
	return nil
}

// Edit 编辑, 需要处理为可按需更新
func (common *CommonService) EditApi(c echo.Context) error {
	// 编辑时必须要传递id
	id := c.Param("id")

	if data, err := common.DealParam(c, common.AddRules); err != nil {
		return err
	} else {
		return common.updateById(id, data)
	}
}
