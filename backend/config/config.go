package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var Config *config

type (
	config struct {
		Server   *server
		PostGres *postgres
	}

	server struct {
		Port string
		Host string
	}

	postgres struct {
		Host     string
		User     string
		Passwrod string
		Dbname   string
		Port     string
	}
)

// 初始化配置文件
func InitConfig(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("config file not found: %s", path)
	}

	Config = new(config)
	_, err := toml.DecodeFile(path, Config)
	if err != nil {
		return fmt.Errorf("decode file: %s", err)
	}

	// fmt.Printf("%#v\n", Config.PostGres)

	return nil
}
