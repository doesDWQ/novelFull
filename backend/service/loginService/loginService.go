package loginService

import (
	"fmt"
	"net/http"
	"time"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com.doesDWQ.novelFull/service/commonService"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type loginService struct {
	commonService.Service
}

func NewLoginService(g *echo.Group) commonService.ServiceInterface {
	service := loginService{}
	service.InitInnerService(&commonService.InnerService{
		Db: db.Db,
		Routes: []*commonService.Route{
			{
				RequestFunc: g.POST,
				Path:        "/adminLogin",
				Func:        service.adminLogin,
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
	})
	return service
}

// adminLogin 后台登录接口
func (l *loginService) adminLogin(c echo.Context) error {
	params, err := l.DealParam(c, map[string]interface{}{
		"user_name": "required",
		"pwd":       "required",
	})

	fmt.Println("params:", params)
	fmt.Println("adminLogin:", err)

	if err != nil {
		fmt.Println("ggggggg:", err)
		return err
	}

	adminUser := model.NewAdminUser()
	searchModel, err := adminUser.FindByUsernamePwd(
		db.Db,
		params["user_name"].(string),
		params["pwd"].(string),
	)
	if err != nil {
		// 未知错误报错
		return err
	}

	auSearch := searchModel.(*model.SearchAdminUserApi)
	// fmt.Printf("auSearch:%#v\n", auSearch.ID)

	// Set custom claims
	claims := &commonService.JwtCustomClaims{
		UserId:   auSearch.ID,
		UserName: auSearch.UserName,
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

	// 更新token
	err = adminUser.UpdateTokenById(db.Db, auSearch.ID, t)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// 退出接口
func (l *loginService) adminLoginOut(c echo.Context) error {
	userId, _ := l.GetTokenInfo(c)
	adminUser := model.NewAdminUser()

	err := adminUser.UpdateTokenById(db.Db, userId, "")
	if err != nil {
		return err
	}

	return nil
}

// 获取当前登录后台用户信息
func (l *loginService) adminUserInfo(c echo.Context) error {
	userId, _ := l.GetTokenInfo(c)

	adminUser := model.NewAdminUser()
	search, err := adminUser.Detail(db.Db, int(userId))
	if err != nil {
		return err
	}

	return l.Success(c, map[string]interface{}{
		"info": search,
	})
}
