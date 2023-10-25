package entity

import "gorm.io/gorm"


type Admin struct {
    gorm.Model
    Username            string                 `json:"username"`
    Email               string                 `json:"email"`
    Password            string                 `json:"password"`
    IsAdmin             bool    `gorm:"default:false" json:"isAdmin"`
	IsVerified          bool    `gorm:"default:false" json:"is_verified"`
}


type AdminErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
