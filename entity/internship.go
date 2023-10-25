package entity

import (
	"time"

	"gorm.io/gorm"
)

type Internship_Listing struct {
    gorm.Model
    Title                  string                       `json:"title"`
    Description            string                       `json:"description"`
    Quota                  int                          `json:"quota"`
    ApplicationForms       []Internship_ApplicationForm `gorm:"foreignKey:InternshipListingID" json:"application_forms"`
    InternshipApplications []Internship_ApplicationForm `gorm:"foreignKey:InternshipListingID" json:"internship_applications"`
    SelectedCandidates     []Selected_Candidate         `gorm:"many2many:selected_candidates" json:"selected_candidates"`
    StatusPendaftaran      string                       `json:"status_pendaftaran"`
    StartDate              time.Time                    `json:"start_date"`
    EndDate                time.Time                    `json:"end_date"`
}

type Selected_Candidate struct {
	UserId 		int
	CandidateID  int
	InternshipID uint
}

type Internship_ApplicationForm struct {
	UserID 				int
	AdminApproval 		bool
	AdminUserID 		int 
	InternshipListingID int                  `json:"internship_listing_id"`
	ApplicationStatusID uint                 `json:"application_status_id"`
	ApplicationStatus   Application_Status   `gorm:"foreignKey:ApplicationStatusID" json:"application_status"`
	// SelectedListings    []Internship_Listing `gorm:"many2many:application_form_selected_listings" json:"selected_listings"`
	// AdminID             uint                 `json:"admin_id" form:"admin_id"`
	CV                  string               `json:"cv" form:"cv"`
	Status              string               `json:"status" form:"status"`
	FirstName           string               `json:"first_name" form:"first_name"`
	LastName            string               `json:"last_name" form:"last_name"`
	Email               string               `json:"email" form:"email"`
	Gender              string               `json:"gender" form:"gender"`
	PhoneNumber         string               `json:"phone_number" form:"phone_number"`
	Address             string               `json:"address" form:"address"`
	City                string               `json:"city" form:"city"`
	State               string               `json:"state" form:"state"`
	PostalCode          string               `json:"postal_code" form:"postal_code"`
	DateOfBirth         string               `json:"date_of_birth" form:"date_of_birth"`
	UniversityOrigin    string               `json:"university_origin" form:"university_origin"`
	UniversityAddress   string               `json:"university_address" form:"university_address"`
	NIM                 string               `json:"nim" form:"nim"`
	GPA                 float64              `json:"gpa" form:"gpa"`
	EducationLevel      string               `json:"education_level" form:"education_level"`
}

type Application_Status struct {
	gorm.Model
	Status string `json:"status" form:"status"`
}
