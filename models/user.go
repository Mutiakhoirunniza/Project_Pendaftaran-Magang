package models

import (
    "github.com/jinzhu/gorm"
)

type User struct {
    gorm.Model
    Username      string `json:"username" form:"username"`
    Fullname      string `json:"fullname" form:"fullname"`
    Email         string `json:"email" form:"email"`
    Password      string `json:"password" form:"password"`
    Gender        string `json:"gender" form:"gender"`
    DateOfBirth   string `json:"date_of_birth" form:"date_of_birth"`
    PhoneNumber   string `json:"phone_number" form:"phone_number"`
    NIM           string `json:"nim" form:"nim"`
    University    string `json:"university" form:"university"`
    UniversityAddress string `json:"university_address" form:"university_address"`
    Major         string `json:"major" form:"major"`
    Semester      string `json:"semester" form:"semester"`
    CVPath        string `json:"cv_path" form:"cv_path"` // Path to the uploaded CV file
    PhotoPath     string `json:"photo_path" form:"photo_path"` // Path to the uploaded photo
}
