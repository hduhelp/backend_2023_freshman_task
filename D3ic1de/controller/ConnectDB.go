package controller

import (
	"TodoList/config"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
)

// 连接数据库
func ConnectToDatabase() (*sql.DB, error) {
	config.InitConfig()
	ms := Database{
		Username: viper.GetString("database.mysql.username"),
		Password: viper.GetString("database.mysql.password"),
		Host:     viper.GetString("database.mysql.host"),
		Port:     viper.GetInt("database.mysql.port"),
		DBname:   viper.GetString("database.mysql.dbname"),
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", ms.Username, ms.Password, ms.Host, ms.Port, ms.DBname)
	db, err := sql.Open("mysql", dataSourceName) // 打开mysql数据库
	if err != nil {
		return nil, err
	}
	err = db.Ping() // 检查连接是否建立，以确保连接存活
	if err != nil {
		return nil, err
	}

	return db, nil
}
