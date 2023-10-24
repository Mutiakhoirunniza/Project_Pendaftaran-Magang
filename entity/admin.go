package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}


type AdminResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
