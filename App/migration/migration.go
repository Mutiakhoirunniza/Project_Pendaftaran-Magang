package App

import (
	"project/models"

	"gorm.io/gorm"
)

func InitMigrationMysql(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Internship{})
	db.AutoMigrate(&models.ApplicationStatus{})
	// Tambahkan migrasi untuk tabel lain jika diperlukan

}
