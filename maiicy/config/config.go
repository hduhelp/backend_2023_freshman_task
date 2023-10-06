package config

import "github.com/pelletier/go-toml"

type Config struct {
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

	return config, nil
}
