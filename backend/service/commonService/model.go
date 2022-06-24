package commonService

import (
	"errors"
	"fmt"
	"time"

	"github.com.doesDWQ.novelFull/db"

	"gorm.io/gorm"
)

// 模型操作基础封装

// add 传输data map类型数据添加数据
func (common *Service) add(data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("添加数据必须传递")
	}

	if data["created_at"] == nil {
		data["created_at"] = time.Now()
	}

	if data["updated_at"] == nil {
		data["updated_at"] = time.Now()
	}

	model := common.Model()
	if err := db.Db.Model(&model).Create(data).Error; err != nil {
		return fmt.Errorf("add error msg: %s", err)
	}
	return nil
}

// detail 根据id获取详情
func (common *Service) detail(id int) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("详情id必须传递")
	}

	model := common.Model()
	searchModel := common.SearchApiModel()
	err := db.Db.Debug().Model(&model).Where("id=?", id).Find(&searchModel).Error
	if err != nil {
		return nil, err
	}
	return &searchModel, nil
}

// updateById 根据id和修改数据map编辑数据
func (common *Service) updateById(id int, data map[string]interface{}) error {
	if id == 0 {
		return fmt.Errorf("update id 必须传递")
	}
	if data == nil {
		return fmt.Errorf("update data is nil")
	}

	if data["updated_at"] == nil {
		data["updated_at"] = time.Now()
	}

	model := common.Model()
	if err := db.Db.Debug().
		Model(&model).
		Where("id=?", id).
		Updates(data).
		Error; err != nil {

		return fmt.Errorf("模型编辑错误：%s", err)
	}
	return nil
}

// deleteById 根据id删除
func (common *Service) deleteById(id int) error {
	if id == 0 {
		return errors.New("delete id 必须传递")
	}

	model := common.Model()
	err := db.Db.Delete(&model, id).Error
	if err != nil {
		return fmt.Errorf("模型删除错误：%s", err)
	}

	return nil
}

// 列表查询, 返回list数据，count条数，查询闭包
func (common *Service) getListQuery(pageInfo *pageInfo) (listModel interface{}, count *int64, query *gorm.DB, err error) {
	if pageInfo == nil {
		err = fmt.Errorf("update data is nil")
		return
	}

	listModel = common.ListModel()
	count = new(int64)

	model := common.Model()
	query = db.Db.Model(&model).
		Debug().
		Count(count).
		Find(&listModel).
		Scopes(common.Paginate(pageInfo.Page, pageInfo.PageSize))

	return
}
