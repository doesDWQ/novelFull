package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminUserService struct{}

// 新增/编辑用户
func (a *AdminUserService) Add(c echo.Context) error {
	user := &model.AdminUser{}
	user.UserName = c.FormValue("username")
	user.Passwrod = c.FormValue("password")
	user.Email = c.FormValue("email")

	fmt.Println(c.Param("id"))
	if c.Param("id") == "" {
		// 新增
		err := db.Db.Create(user).Error
		if err != nil {
			return err
		}
	} else {
		// 编辑
		id := c.Param("id")
		i, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		user.Model.ID = uint(i)
		err = db.Db.Debug().Save(user).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除用户
func (a *AdminUserService) Delete(c echo.Context) error {
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

// 用户列表
func (a *AdminUserService) List(c echo.Context) error {
	pageSize := c.FormValue("limit")
	page := c.FormValue("page")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return err
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	users := new([]model.AdminUser)
	var count int64
	err = db.Db.Find(users).Count(&count).Limit(pageSizeInt).Offset(pageSizeInt*pageInt - 1).Error
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

// 用户详情
func (a *AdminUserService) Detail(c echo.Context) error {

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
	err = db.Db.First(user).Error
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
