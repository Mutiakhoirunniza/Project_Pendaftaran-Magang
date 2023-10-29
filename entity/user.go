package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string                       `json:"username" form:"username" gorm:"unique;not null"`
	Password          string                       `json:"password" form:"password" gorm:"not null"`
	Email             string                       `json:"email" form:"email" gorm:"unique;not null"`
	Gender            string                       `json:"gender" form:"gender"`
	PhoneNumber       string                       `json:"phone_number" form:"phone_number"`
	UniversityName    string                       `json:"university_name" form:"university_name"`
	UniversityAddress string                       `json:"university_address" form:"university_address"`
	Major             string                       `json:"major" form:"major"`
	Form              []Internship_ApplicationForm //`gorm:"foreignKey:UserID"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	Token string `json:"token"`
}
