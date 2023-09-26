package models

import (
	"bytes"
	_ "embed"
	"github.com/spf13/viper"
	"os"
)

var Env *viper.Viper

//go:embed default.yml
var defaultConf []byte

func InitConfig() {
	Env = viper.New()
	Env.SetConfigType("yaml")
	Env.SetConfigName("config")
	Env.AddConfigPath(".")
	if err := Env.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := os.WriteFile("./config.yml", defaultConf, 0666); err != nil {
				Logger.Fatalln(err)
			}
			_ = Env.ReadConfig(bytes.NewReader(defaultConf))
		} else {
			Logger.Fatalln(err)
		}
	}
}
