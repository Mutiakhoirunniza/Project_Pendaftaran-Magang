package App

import (
	"fmt"
	App "project/App/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB // Objek GORM untuk koneksi database

func InitDBMysql(AppConfig *App.AppConfig) *gorm.DB {

	// declare struct config & variable connectionString
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUSER, AppConfig.DBPASS, AppConfig.DBHOST, AppConfig.DBPORT, AppConfig.DBNAME)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return db
}
