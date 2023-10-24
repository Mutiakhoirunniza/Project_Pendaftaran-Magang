package migration

import (
    "miniproject/entity"
    "gorm.io/gorm"
)

func InitMigrationMysql(db *gorm.DB) {
    db.AutoMigrate(
        &entity.Admin{},
        &entity.User{},
        &entity.InternshipListing{},
        &entity.InternshipApplicationForm{},
        &entity.ApplicationStatus{},
    )
}

