package models

import (
	"gorm.io/gorm"
)

type Internship struct {
    gorm.Model
    Title       string `gorm:"not null"`
    Description string
    Quota  int
}
