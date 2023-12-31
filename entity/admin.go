package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `json:"username" form:"username" gorm:"unique;not null"`
	Email    string `json:"email" form:"email" gorm:"unique;not null"`
	Password string `json:"password" form:"password" gorm:"not null"`
	Role     string `gorm:"type:enum('user','admin');default:'admin'"`
}

type AdminResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
