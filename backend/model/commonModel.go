package model

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com.doesDWQ.novelFull/types"
	"gorm.io/gorm"
)

type CommonModelInterface interface {
	Add(tx *gorm.DB, data map[string]interface{}) error
	Detail(tx *gorm.DB, id int) (search interface{}, err error)
	UpdateById(tx *gorm.DB, id int, data map[string]interface{}) error
	DeleteById(tx *gorm.DB, id int) error
	GetListQuery(tx *gorm.DB, pageInfo *types.PageInfo) (listModel interface{}, count *int64, query *gorm.DB, err error)
}

type CommonModel struct {
	// 持有当前model的引用
	currentModel interface{}
	// 单个查询model获取函数
	modelFunc func() interface{}
	// 列表查询获取函数
	listModelFunc func() interface{}
}

func (common *CommonModel) SetModelFunc(mdelFunc func() interface{}) {
	common.modelFunc = mdelFunc
}

func (common *CommonModel) GetModelFunc() (searchModel interface{}, err error) {
	if common.modelFunc == nil {
		err = errors.New("mdelFunc not set")
		return
	}
	searchModel = common.modelFunc()
	return
}

func (common *CommonModel) SetListModelFunc(listModelFunc func() interface{}) {
	common.listModelFunc = listModelFunc
}

func (common *CommonModel) GetListModel() (listModel interface{}, err error) {
	if common.modelFunc == nil {
		err = errors.New("listModelFunc not set")
		return
	}
	listModel = common.listModelFunc()
	return
}

func (common *CommonModel) SetCurrentModel(currentModel interface{}) {
	common.currentModel = currentModel
}

// 主要用来判定 currentModel 是否设置
func (common *CommonModel) GetCurrentModel() (currentModel interface{}, err error) {
	if common.currentModel == nil {
		err = errors.New("not set currentModel")
		return
	}

	currentModel = common.currentModel
	return
}

// add 传输data map类型数据添加数据
func (common CommonModel) Add(tx *gorm.DB, data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("add data is empty")
	}

	if data["created_at"] == nil {
		data["created_at"] = time.Now()
	}

	if data["updated_at"] == nil {
		data["updated_at"] = time.Now()
	}

	currentModel, err := common.GetCurrentModel()
	if err != nil {
		return err
	}
	if err := tx.Debug().Model(currentModel).Create(data).Error; err != nil {
		return fmt.Errorf("add error msg: %s", err)
	}
	return nil
}

// detail 根据id获取详情
func (common CommonModel) Detail(tx *gorm.DB, id int) (search interface{}, err error) {
	fmt.Println("后台用户id：", id)
	if id == 0 {
		return nil, fmt.Errorf("详情id必须传递")
	}

	search, err = common.GetModelFunc()
	if err != nil {
		return
	}

	currentModel, err := common.GetCurrentModel()
	if err != nil {
		return nil, err
	}
	err = tx.Debug().Model(currentModel).Where("id=?", id).Find(search).Error
	if err != nil {
		return nil, fmt.Errorf("detail error: %s", err.Error())
	}

	// 判断数据是否存在
	// 反射查询id值
	rfValue := reflect.ValueOf(search).Elem()
	reflectId := rfValue.FieldByName("ID").Uint()
	if int(reflectId) != id {
		return nil, errors.New("record not found")
	}

	return search, nil
}

// 判断记录是否存在
func (common *CommonModel) RecordIsExists(tx *gorm.DB, id int) error {
	var search interface{}
	search, err := common.Detail(tx, id)
	if err != nil {
		return err
	}

	// 反射查询id值
	rfValue := reflect.ValueOf(search).Elem()
	reflectId := rfValue.FieldByName("ID").Uint()
	if int(reflectId) != id {
		return errors.New("record not found")
	}

	return nil
}

// updateById 根据id和修改数据map编辑数据
func (common CommonModel) UpdateById(tx *gorm.DB, id int, data map[string]interface{}) error {
	if id == 0 {
		return fmt.Errorf("update id not found")
	}
	if data == nil {
		return fmt.Errorf("update data is nil")
	}
	// 判断需要修改的数据是否存在
	err := common.RecordIsExists(tx, id)
	if err != nil {
		return err
	}

	if data["updated_at"] == nil {
		data["updated_at"] = time.Now()
	}

	currentModel, err := common.GetCurrentModel()
	if err != nil {
		return err
	}

	if err := tx.Debug().
		Model(currentModel).
		Where("id=?", id).
		Updates(data).
		Error; err != nil {

		return fmt.Errorf("模型编辑错误：%s", err)
	}
	return nil
}

// deleteById 根据id删除
func (common CommonModel) DeleteById(tx *gorm.DB, id int) error {
	if id == 0 {
		return errors.New("delete id 必须传递")
	}

	currentModel, err := common.GetCurrentModel()
	if err != nil {
		return err
	}

	// 判断数据是否存在
	err = common.RecordIsExists(tx, id)
	if err != nil {
		return err
	}

	var search interface{}

	// 反射查询id值
	rfValue := reflect.ValueOf(search).Elem()
	reflectId := rfValue.FieldByName("ID").Uint()
	if int(reflectId) != id {
		return errors.New("delete record is not found")
	}

	err = tx.Debug().Delete(currentModel, id).Error
	if err != nil {
		return fmt.Errorf("mode delete error：%s", err)
	}

	return nil
}

// 列表查询, 返回list数据，count条数，查询闭包
func (common CommonModel) GetListQuery(tx *gorm.DB, pageInfo *types.PageInfo) (listModel interface{}, count *int64, query *gorm.DB, err error) {
	if pageInfo == nil {
		err = fmt.Errorf("update data is nil")
		return
	}

	listModel, err = common.GetListModel()
	if err != nil {
		return
	}
	count = new(int64)

	var currentModel interface{}
	currentModel, err = common.GetCurrentModel()
	if err != nil {
		return
	}

	query = tx.Model(currentModel).
		Debug().
		Count(count).
		Find(listModel).
		Scopes(common.Paginate(pageInfo.Page, pageInfo.PageSize))

	return
}

// Paginate 分页
func (model *CommonModel) Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fmt.Println("paginate")
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
