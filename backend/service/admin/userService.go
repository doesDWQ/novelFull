package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserService struct {
	service.CommonService
}

// Add 新增用户
func (a *UserService) Add(c echo.Context) error {
	// gorm.Model 的定义
	type userRequest struct {
		UserName string `validate:"required"`
		Passwrod string `validate:"required"`
		Email    string `validate:"required"`
	}

	request := &userRequest{}
	err := c.Bind(request)
	if err != nil {
		return fmt.Errorf("bind err: %s", err)
	}

	fmt.Printf("%#v\n", request)

	user := &model.AdminUser{}
	err = copier.CopyWithOption(user, request, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return fmt.Errorf("copy err: %s", err)
	}

	// 新增
	err = db.Db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// Edit 编辑用户
func (a *UserService) Edit(c echo.Context) error {
	user := &model.AdminUser{}
	user.UserName = c.FormValue("username")
	user.Passwrod = c.FormValue("password")
	user.Email = c.FormValue("email")

	id := c.Param("id")
	if id == "" {
		return errors.New("id不能为空！")
	}

	// 编辑
	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	user.Model.ID = uint(i)
	err = db.Db.Debug().Save(user).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete 删除用户
func (a *UserService) Delete(c echo.Context) error {
	id := c.Param("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	user := &model.AdminUser{
		Model: gorm.Model{
			ID: uint(i),
		},
	}

	err = db.Db.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// List 用户列表
func (a *UserService) List(c echo.Context) error {
	users := new([]model.AdminUserApi)
	var count int64
	err := db.Db.Model(&model.AdminUser{}).Scan(users).Count(&count).Scopes(a.Paginate(c.Request())).Error
	if err != nil {
		return err
	}

	err = c.JSON(http.StatusOK, echo.Map{
		"data": map[string]interface{}{
			"users": users,
			"cnt":   count,
		},
	})

	if err != nil {
		return err
	}
	return nil
}

// Detail 用户详情
func (a *UserService) Detail(c echo.Context) error {

	id := c.Param("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	user := &model.AdminUserApi{
		Model: gorm.Model{
			ID: uint(i),
		},
	}
	err = db.Db.Model(&model.AdminUser{}).Scan(user).Where(user).First(user).Error
	if err != nil {
		return err
	}

	err = c.JSON(http.StatusOK, echo.Map{
		"data": map[string]interface{}{
			"user": user,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
