package database

import (
	"fmt"
	"miniproject/infra/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDBMysql(cfg *config.AppConfig) {
    var err error
    // Format string koneksi MySQL
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.DBUSER, cfg.DBPASS, cfg.DBHOST, cfg.DBPORT, cfg.DBNAME)
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Gagal menghubungkan ke database MySQL")
    }
}

// var db *gorm.DB // Objek GORM untuk koneksi database

// func InitDBMysql(cfg *config.AppConfig) *gorm.DB {
// 	var err error
// 	// declare struct config & variable connectionString
// 	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		cfg.DBUSER, cfg.DBPASS, cfg.DBHOST, cfg.DBPORT, cfg.DBNAME)

// 	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }
