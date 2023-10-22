package migration

import (
	"miniproject/entity"

	"gorm.io/gorm"
)

func InitMigrationMysql(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Admin{})

	// db.AutoMigrate(&entity.Internship{})
	// db.AutoMigrate(&entity.ApplicationStatus{})
	// Tambahkan migrasi 
}
