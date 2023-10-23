package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID            int    `gorm:"primaryKey" form:"user_id"`
	Username          string `json:"username" form:"username" gorm:"unique;not null"`
	Password          string `json:"password" form:"password" gorm:"not null"`
	Email             string `json:"email" form:"email" gorm:"unique;not null"`
	Gender            string `json:"gender" form:"gender"`
	PhoneNumber       string `json:"phone_number" form:"phone_number"`
	UniversityName    string `json:"university_name" form:"university_name"`
	UniversityAddress string `json:"university_address" form:"university_address"`
	Major             string `json:"major" form:"major"`
	ProfilePicture    string `json:"profile_picture" form:"profile_picture"`
}
