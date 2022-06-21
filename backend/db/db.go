package db

import (
	"fmt"

	"github.com.doesDWQ.novelFull/config"
	"github.com.doesDWQ.novelFull/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

// 初始化数据库
func InitDb() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.Config.PostGres.Host,
		config.Config.PostGres.User,
		config.Config.PostGres.Passwrod,
		config.Config.PostGres.Dbname,
		config.Config.PostGres.Port,
	)
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("init db: %s", err)
	}

	migrate()
	return nil
}

func migrate() {
	Db.AutoMigrate(
		&model.AdminUser{},
	)
}
