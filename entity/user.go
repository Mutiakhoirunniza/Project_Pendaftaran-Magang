package entity

import "gorm.io/gorm"

type User struct {
    gorm.Model
    UserID            int    `gorm:"primaryKey"`
    Username          string `json:"username" gorm:"unique;not null"`
    Password          string `json:"password" gorm:"not null"`
    Email             string `json:"email" gorm:"unique;not null"`
    Gender            string `json:"gender"`
    PhoneNumber       string `json:"phone_number"`
    UniversityName    string `json:"university_name"`
    UniversityAddress string `json:"university_address"`
    Major             string `json:"major"`
    ProfilePicture    string `json:"profile_picture"`
}
