// app/models/application_status.go
package models

import "gorm.io/gorm"

type ApplicationStatus struct {
    gorm.Model
    UserID      uint   `gorm:"not null"`
    InternshipID uint   `gorm:"not null"`
    Status      string `gorm:"not null;check:Status IN ('Pengajuan', 'Diverifikasi = 'diterima, ditolak', 'Dibatalkan')"`
    
}

