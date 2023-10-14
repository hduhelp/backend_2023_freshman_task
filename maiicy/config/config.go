package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
)

type Config struct {
	DBPath     string `toml:"db_path"`
	SecretKey  string `toml:"secret_key"`
	ServerIP   string `toml:"server_ip"`
	ServerPort int    `toml:"server_port"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	tree, err := toml.LoadFile(filename)
	if err != nil {
		return config, err
	}

	if err := tree.Unmarshal(&config); err != nil {
		return config, err
	}

	// 从环境变量中获取 SecretKey
	secretKeyEnv := os.Getenv("SECRET_KEY")
	if secretKeyEnv == "" {
		fmt.Println("错误：未设置服务器的 JWT 秘钥 (SECRET_KEY)。")
		fmt.Println("请设置环境变量 SECRET_KEY 以确保服务器的安全性。")
		os.Exit(1) // 退出应用程序，因为缺少必要的秘钥
	} else {
		config.SecretKey = secretKeyEnv
	}

	return config, nil
}
