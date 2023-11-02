package entity

import (
	"gorm.io/gorm"
)

type Internship_Listing struct {
	gorm.Model
	Title            string                       `json:"title" form:"title"`
	Description      string                       `json:"description" form:"description"`
	Quota            int                          `json:"quota" form:"quota"`
	Qualifications   string                       `json:"qualifications" form:"qualifications"`
	StartDate        string                       `json:"start_date" form:"start_date"`
	EndDate          string                       `json:"end_date" form:"end_date"`
	Qualifications   string                       `json:"qualifications" form:"qualifications"`
	ApplicationForms []Internship_ApplicationForm `gorm:"foreignKey:InternshipListingID" json:"applicationforms" form:"applicationforms"`
}

type Internship_ApplicationForm struct {
	gorm.Model
	CV                  string               `json:"cv" form:"cv"`
	Nim                 string               `json:"nim" form:"nim"`
	GPA                 float64              `json:"gpa" form:"gpa"`
	EducationLevel      string               `json:"education_level" form:"education_level"`
	UserID              int                  `json:"UserID " form:"UserID"`
	Status              string               `json:"status"`
	UserEmail           string               `json:"user_email" form:"user_email" gorm:"not null"`
	Username            string               `json:"username" form:"username" gorm:"not null"`
	SelectedTitle       string               `json:"selected_title" form:"selected_title"`
	IsCanceled          bool                 `json:"is_canceled" form:"is_canceled"`
	InternshipListingID uint                 `json:"internshiplistingID" form:"internshiplistingID" gorm:"not null"`
	Selected_Candidates []Selected_Candidate `gorm:"foreignKey:InternshipApplicationFormID" json:"selected_candidates" form:"selected_candidates"`
}

type Selected_Candidate struct {
	gorm.Model
	InternshipApplicationFormID uint
	InternshipApplicationForm   Internship_ApplicationForm `gorm:"foreignKey:InternshipApplicationFormID"`
}
