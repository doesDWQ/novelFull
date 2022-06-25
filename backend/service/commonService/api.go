package commonService

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Add 新增
func (common *Service) AddAPi(c echo.Context) error {
	data, err := common.DealParam(c, common.innerService.AddRules)
	if err != nil {
		return err
	} else {
		if err != nil {
			return err
		}
		return common.innerService.Model().Add(common.innerService.Db, data)
	}
}

// Detail 详情
func (common *Service) DetailApi(c echo.Context) error {

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("detail strconv error: %s", err)
	}

	model := common.innerService.Model()
	search, err := model.Detail(common.innerService.Db, idInt)
	if err != nil {
		return common.Error(c, err)
	} else {
		return common.Success(c, map[string]interface{}{
			"detail": search,
		})
	}
}

// Delete 删除
func (common *Service) DeleteApi(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	model := common.innerService.Model()
	return common.Error(c, model.DeleteById(common.innerService.Db, idInt))
}

// List 列表
func (common *Service) ListApi(c echo.Context) error {
	// 获取分页信息
	pageInfo, err := common.GetPageInfo(c)
	if err != nil {
		return err
	}

	model := common.innerService.Model()
	listData, count, query, err := model.GetListQuery(common.innerService.Db, pageInfo)
	if err != nil {
		return err
	}

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
func (common *Service) EditApi(c echo.Context) error {
	// 编辑时必须要传递id
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return errors.New("id convert failed")
	}

	data, err := common.DealParam(c, common.innerService.AddRules)
	if err != nil {
		return err
	} else {
		model := common.innerService.Model()
		return common.Error(c, model.UpdateById(common.innerService.Db, id, data))
	}
}
