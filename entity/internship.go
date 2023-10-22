package entity

import (
	"time"

	"gorm.io/gorm"
)

type InputData struct {
	gorm.Model
    Title       string `json:"title"`
	Description string `json:"description"`
	Quota       int    `json:"quota"`
	CreatedDate string `json:"created_date"`
}

type InternshipListing struct {
	gorm.Model
    ID                 uint      `gorm:"primaryKey" json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Quota              int       `json:"quota"`
	CreatedDate        time.Time `json:"created_date"`
	SelectedCandidates []int     `json:"selected_candidates"`
	StatusPendaftaran string	`json:"status_pendaftaran"`
}

type InternshipApplicationForm struct {
	gorm.Model
    ID 					int 
	Status  			string
	Fullname            string    `json:"fullname"`
	NIM                 string    `json:"nim"`
	JurusanProdi        string    `json:"jurusan_prodi"`
	PhoneNumber         string    `json:"phone_number"`
	Gender              string    `json:"gender"`
	DateOfBirth         time.Time `json:"date_of_birth"`
	UniversityOrigin    string    `json:"university_origin"`
	UniversityAddress   string    `json:"university_address"`
	GPA                 float64   `json:"gpa"`
	OrganizationalExp   string    `json:"organizational_experience"`
	InternshipStartDate time.Time `json:"internship_start_date"`
	InternshipEndDate   time.Time `json:"internship_end_date"`
	CV                  string    `json:"cv"`
	InternshipListingID uint      `json:"internship_listing_id"`
	Email               string    `json:"email"`
	Address             string    `json:"address"`
}

type ApplicationStatus struct {
	gorm.Model
	InternshipApplicationFormID uint
	Status                      string
}
