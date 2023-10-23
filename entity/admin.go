package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
}
