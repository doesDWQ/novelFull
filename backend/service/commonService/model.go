package commonService

import (
	"fmt"

	"github.com.doesDWQ.novelFull/db"

	"gorm.io/gorm"
)

// 模型操作基础封装

// add 传输data map类型数据添加数据
func (common *CommonService) add(data map[string]interface{}) error {
	if err := db.Db.Model(common.Model()).Create(data).Error; err != nil {
		return fmt.Errorf("add error msg: %s", err)
	}
	return nil
}

// detail 根据id获取详情
func (common *CommonService) detail(id int) (interface{}, error) {
	searchModel := common.SearchApiModel()
	err := db.Db.Model(common.Model()).Scan(searchModel).First(searchModel, id).Error
	if err != nil {
		return nil, err
	}
	return err, nil
}

// updateById 根据id和修改数据map编辑数据
func (common *CommonService) updateById(id string, data map[string]interface{}) error {
	if err := db.Db.Debug().
		Model(common.Model()).
		Where("id=?", id).
		Updates(data).
		Error; err != nil {

		return err
	}
	return nil
}

// deleteById 根据id删除
func (common *CommonService) deleteById(id int) error {
	err := db.Db.Delete(common.Model(), id).Error
	if err != nil {
		return err
	}

	return nil
}

// 列表查询, 返回list数据，count条数，查询闭包
func (common *CommonService) getListQuery(pageInfo *pageInfo) (listData interface{}, count *int64, query *gorm.DB) {
	listData = common.ListModel()
	count = new(int64)

	query = db.Db.Model(common.Model()).
		Scan(listData).
		Count(count).
		Scopes(common.Paginate(pageInfo.Page, pageInfo.PageSize))

	return
}
