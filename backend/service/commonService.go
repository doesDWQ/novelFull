package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CommonService struct {
	// 获取原始model的方法
	Model func() interface{}
	// 列表model
	ListModel func() interface{}
	// 查询model
	SearchApiModel func() interface{}
	// 添加时的请求参数model
	AddRequestModel func() interface{}
	// 编辑时的请求参数model
	EditRequestModel func() interface{}
}

type AdminJwtCustomClaims struct {
	UserId   uint   `json:user_id`
	UserName string `json:"user_name"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

// Paginate 分页
func (common *CommonService) Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.Form
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
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

// Add 新增
func (common *CommonService) Add(c echo.Context) error {
	request := common.AddRequestModel()
	err := c.Bind(request)
	if err != nil {
		return fmt.Errorf("bind err: %s", err)
	}

	err = c.Validate(request)
	if err != nil {
		return common.Error(c, fmt.Sprintf("validate err: %s", err))
	}

	fmt.Printf("%#v\n", request)

	model := common.Model()
	err = copier.CopyWithOption(model, request, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return fmt.Errorf("copy err: %s", err)
	}

	// 新增
	err = db.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Edit 编辑, 需要处理为可按需更新
func (common *CommonService) Edit(c echo.Context) error {
	request := common.EditRequestModel()
	err := c.Bind(request)
	if err != nil {
		return fmt.Errorf("bind error: %s", err)
	}

	err = c.Validate(request)
	if err != nil {
		return common.Error(c, fmt.Sprintf("validate err: %s", err))
	}

	model := common.Model()
	err = copier.CopyWithOption(model, request, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return fmt.Errorf("copy err: %s", err)
	}
	id := c.Param("id")
	if id == "" {
		return errors.New("id不能为空！")
	}

	// 编辑 需要重新编写这一块逻辑，使用map
	// i, err := strconv.Atoi(id)
	// if err != nil {
	// 	return err
	// }

	// m := structs.Map(model)
	// m["id"] = uint(i)
	// err = db.Db.Debug().Model(model).Updates(m).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Delete 删除
func (common *CommonService) Delete(c echo.Context) error {
	id := c.Param("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	err = db.Db.Delete(common.Model(), uint(i)).Error
	if err != nil {
		return err
	}

	return nil
}

// List 列表
func (common *CommonService) List(c echo.Context) error {
	listModel := common.ListModel()
	var count int64
	err := db.Db.Model(common.Model()).
		Scan(listModel).
		Count(&count).
		Scopes(common.Paginate(c.Request())).
		Error
	if err != nil {
		return err
	}

	err = common.Success(c, map[string]interface{}{
		"users": listModel,
		"cnt":   count,
	})

	if err != nil {
		return err
	}
	return nil
}

// Detail 详情
func (common *CommonService) Detail(c echo.Context) error {

	id := c.Param("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("detail strconv error: %s", err)
	}
	searchModel := common.SearchApiModel()
	err = db.Db.Model(common.Model()).Scan(searchModel).First(searchModel, uint(i)).Error
	if err != nil {
		return err
	}

	err = common.Success(c, map[string]interface{}{
		"user": searchModel,
	})

	if err != nil {
		return err
	}

	return nil
}

// 返回结构
func (common *CommonService) Success(c echo.Context, data map[string]interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"data":   data,
		"status": "1",
		"msg":    "",
	})
}

// 返回结构
func (common *CommonService) Error(c echo.Context, errMsg string) error {
	return c.JSON(http.StatusOK, echo.Map{
		"data":   nil,
		"status": "0",
		"msg":    errMsg,
	})
}
