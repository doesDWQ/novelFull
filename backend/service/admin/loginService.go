package admin

import (
	"errors"
	"net/http"
	"time"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type loginService struct {
	service.CommonService
}

func NewLoginService() *loginService {
	return &loginService{}
}

func (l *loginService) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 查询用户名和密码是否匹配
	user := &model.AdminUser{
		UserName: username,
		Passwrod: password,
	}

	// Throws unauthorized error
	err := db.Db.Debug().Where(user).First(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未查询到, 说明用户不存在
			return l.Error(c, "登录错误")
		} else {
			// 未知错误报错
			return err
		}
	}

	// Set custom claims
	claims := &service.AdminJwtCustomClaims{
		UserId:   user.ID,
		UserName: user.UserName,
		Admin:    true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.Config.Jwt.Secret))
	if err != nil {
		return err
	}

	user.Token = t
	// 更新token
	err = db.Db.Save(user).Error
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// 校验token是否存在
func (l *loginService) VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if l.SkipVerify(c) {
			// 过权限校验
			return next(c)
		}

		// 校验token
		cs := service.CommonService{}
		userId, token := cs.GetTokenInfo(c)
		user := &model.AdminUser{
			Model: gorm.Model{
				ID: userId,
			},
		}

		err := db.Db.Where(user).First(user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return l.Error(c, "用户id不存在")
			}
			return err
		}

		// token为空表示未登录，token和表里面的token没有对照上也是登录错误
		if user.Token == "" && user.Token != token.Raw {
			// token校验不存在则报错
			return l.Error(c, "用户信息校验失败")
		}

		return next(c)
	}
}

// 跳过权限校验的路由
func (l *loginService) SkipVerify(c echo.Context) bool {
	skipApis := map[string]struct{}{
		"/admin/login": {},
	}

	url := c.Request().URL.String()

	if _, exists := skipApis[url]; exists {
		return true
	}

	return false
}

// 退出接口
func (l *loginService) LoginOut(c echo.Context) error {
	cs := service.CommonService{}
	userId, _ := cs.GetTokenInfo(c)

	user := &model.AdminUser{Model: gorm.Model{
		ID: userId,
	}}

	err := db.Db.Model(&user).Update("token", "").Error
	if err != nil {
		return err
	}

	return nil
}

// 获取当前登录后台用户信息
func (l *loginService) AdminUserInfo(c echo.Context) error {
	userId, _ := l.GetTokenInfo(c)

	user := &model.AdminUser{
		Model: gorm.Model{ID: userId},
	}

	userInfo := &model.AdminUserApi{}
	err := db.Db.Debug().Model(user).Where(user).Scan(userInfo).Error
	if err != nil {
		return err
	}

	return l.Success(c, map[string]interface{}{
		"info": userInfo,
	})
}
