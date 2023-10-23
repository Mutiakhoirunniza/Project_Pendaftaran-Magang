package entity

import (
	"time"

	"gorm.io/gorm"
)

type InputData struct {
	gorm.Model
	Title       string    `json:"title" form:"title"`
	Description string    `json:"description" form:"description"`
	Quota       int       `json:"quota" form:"quota"`
	CreatedDate time.Time `json:"created_date" form:"created_date"`
}

type InternshipListing struct {
	gorm.Model
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Title              string    `json:"title" form:"title"`
	Description        string    `json:"description" form:"description"`
	Quota              int       `json:"quota" form:"quota"`
	CreatedDate        time.Time `json:"created_date" form:"created_date"`
	SelectedCandidates []int     `json:"selected_candidates" form:"selected_candidates"`
	StatusPendaftaran  string    `json:"status_pendaftaran" form:"status_pendaftaran"`
}

type InternshipApplicationForm struct {
	gorm.Model
	ID                  int
	Status              string
	Fullname            string    `json:"fullname" form:"fullname"`
	NIM                 string    `json:"nim" form:"nim"`
	JurusanProdi        string    `json:"jurusan_prodi" form:"jurusan_prodi"`
	PhoneNumber         string    `json:"phone_number" form:"phone_number"`
	Gender              string    `json:"gender" form:"gender"`
	DateOfBirth         time.Time `json:"date_of_birth" form:"date_of_birth"`
	UniversityOrigin    string    `json:"university_origin" form:"university_origin"`
	UniversityAddress   string    `json:"university_address" form:"university_address"`
	GPA                 float64   `json:"gpa" form:"gpa"`
	OrganizationalExp   string    `json:"organizational_experience" form:"organizational_experience"`
	InternshipStartDate time.Time `json:"internship_start_date" form:"internship_start_date"`
	InternshipEndDate   time.Time `json:"internship_end_date" form:"internship_end_date"`
	CV                  string    `json:"cv" form:"cv"`
	InternshipListingID uint      `json:"internship_listing_id" form:"internship_listing_id"`
	Email               string    `json:"email" form:"email"`
	Address             string    `json:"address" form:"address"`
}

type ApplicationStatus struct {
	gorm.Model
	InternshipApplicationFormID int    `json:"internship_application_form_id" form:"internship_application_form_id"`
	Status                      string `json:"status" form:"status"`
}
