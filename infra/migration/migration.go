package migration

import (
    "miniproject/entity"
    "gorm.io/gorm"
)

func InitMigrationMysql(db *gorm.DB) {
    db.AutoMigrate(
        &entity.User{},
        &entity.Admin{},
        &entity.Internship_Listing{}, 
        &entity.Internship_ApplicationForm{}, 
        &entity.Selected_Candidate{})
}

