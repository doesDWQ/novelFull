package loginService

import (
	"errors"
	"net/http"
	"time"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service/commonService"
	"github.com.doesDWQ.novelFull/types"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type loginService struct {
	commonService.Service
}

func NewLoginService(g *echo.Group) types.Service {
	service := &loginService{}
	service.Service = commonService.Service{
		Model: func() interface{} {
			return nil
		},
		ListModel: func() interface{} {
			return nil
		},
		SearchApiModel: func() interface{} {
			return nil
		},
		AddRules:  map[string]interface{}{},
		EditRules: map[string]interface{}{},
		Routes: []*types.Route{
			{
				RequestFunc: g.POST,
				Path:        "/adminLogin",
				Func:        service.adminLogin,
				// 登录接口需要跳过权限校验
				SkipVerify: true,
			},
			{
				RequestFunc: g.POST,
				Path:        "/adminLoginout",
				Func:        service.adminLoginOut,
			},
			{
				RequestFunc: g.GET,
				Path:        "/adminUserInfo",
				Func:        service.adminUserInfo,
			},
		},
		Group:        g,
		AutoRegister: false,
	}
	return service
}

// adminLogin 后台登录接口
func (l *loginService) adminLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 查询用户名和密码是否匹配
	user := &model.AdminUser{
		UserName: username,
		Pwd:      password,
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
	claims := &commonService.JwtCustomClaims{
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

// 退出接口
func (l *loginService) adminLoginOut(c echo.Context) error {
	cs := commonService.Service{}
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
func (l *loginService) adminUserInfo(c echo.Context) error {
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
