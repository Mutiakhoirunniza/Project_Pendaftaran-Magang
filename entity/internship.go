package entity

import (
	"time"

	"github.com/jinzhu/gorm"
)

// type InputData struct {
// 	gorm.Model
// 	Title       string    `json:"title" form:"title"`
// 	Description string    `json:"description" form:"description"`
// 	Quota       int       `json:"quota" form:"quota"`
// 	CreatedDate time.Time `json:"created_date" form:"created_date"`
// }

type InternshipListing struct {
	ID         			 	uint      `gorm:"primaryKey" json:"id"`
	Title      			 	string    `json:"title"`
	Description 			string    `json:"description"`
	Quota       			int       `json:"quota"`
	SelectedCandidates []SelectedCandidate `gorm:"foreignKey:InternshipID"`
	StatusPendaftaran 		string 
	CreatedDate 			time.Time `json:"created_date"`
	UpdatedAt         		time.Time `json:"updated_at"`
}


type SelectedCandidate struct {
    ID             uint `gorm:"primaryKey" json:"id"`
    InternshipID   uint
    CandidateID    int
}


type InternshipApplicationForm struct {
	ID               	int   `gorm:"primaryKey" json:"id"`
	InternshipListingID int   `json:"internship_listing_id"`
	CV               	string `json:"cv"`
	Status          	 string `json:"status"`
	FirstName         	string `json:"first_name"`
	LastName          	string `json:"last_name"`
	Email             	string `json:"email"`
	Gender            	string  `json:"gender" form:"gender"`
	PhoneNumber       	string `json:"phone_number"`
	Address           	string `json:"address"`
	City              	string `json:"city"`
	State             	string `json:"state"`
	PostalCode        	string `json:"postal_code"`
	DateOfBirth       	string `json:"date_of_birth"`
	UniversityOrigin  	string `json:"university_origin" form:"university_origin"`
	UniversityAddress 	string  `json:"university_address" form:"university_address"`
	NIM               	string  `json:"nim" form:"nim"`
	GPA               	float64   `json:"gpa" form:"gpa"`
	EducationLevel    	string `json:"education_level"`
	// Add other necessary fields here
}


type ApplicationStatus struct {
	gorm.Model
	ID                        	int   `gorm:"primaryKey" json:"id"`
	InternshipApplicationFormID int    `json:"internship_application_form_id" form:"internship_application_form_id"`
	Status                      string `json:"status" form:"status"`
}


