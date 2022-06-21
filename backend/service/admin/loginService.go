package admin

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com.doesDWQ.novelFull/db"
	"github.com.doesDWQ.novelFull/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LoginService struct{}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type AdminJwtCustomClaims struct {
	ID       uint   `json:id`
	UserName string `json:"user_name"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func (l *LoginService) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 查询用户名和密码是否匹配
	user := &model.AdminUser{
		UserName: username,
		Passwrod: password,
	}

	// Throws unauthorized error
	err := db.Db.Debug().Where(user).First(user).Error

	fmt.Println("err:", err, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未查询到, 说明用户不存在
			return echo.ErrUnauthorized
		} else {
			// 未知错误报错
			return err
		}
	}

	fmt.Println(*user)

	// Set custom claims
	claims := &AdminJwtCustomClaims{
		user.ID,
		user.UserName,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
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
