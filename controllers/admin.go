package controllers

import (
	"net/http"
	"project/models"

	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// InternshipController adalah tipe controller yang menangani operasi lowongan magang
type InternshipController struct {
	DB *gorm.DB
}

// NewInternshipController adalah fungsi untuk membuat instance InternshipController
func NewInternshipController(db *gorm.DB) *InternshipController {
	return &InternshipController{DB: db}
}

// PastikanPenggunaAdmin memeriksa apakah pengguna adalah admin atau bukan.
func (ic *InternshipController) PastikanPenggunaAdmin(c echo.Context) error {
	isAdmin, ok := c.Get("isAdmin").(bool)
	if !ok || !isAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, "Anda tidak memiliki izin admin")
	}
	return nil
}

// AdminLogin digunakan untuk mengotentikasi admin.
func (ic *InternshipController) AdminLogin(c echo.Context) error {
	// Menerima data login admin dari permintaan
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Cari admin berdasarkan username
	var admin models.Admin
	if err := ic.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Login admin gagal"})
	}

	// Verifikasi kata sandi
	if !checkPassword(password, admin.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Login admin gagal"})
	}

	// Setelah berhasil terotentikasi, Anda dapat memberikan token JWT sebagai respons
	token, err := helpers.GenerateJWT(admin.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Terjadi kesalahan saat menghasilkan token JWT"})
	}

	// Return respons dengan token
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// GetApplicationsByStatus digunakan untuk admin melihat daftar aplikasi berdasarkan status.
func (ic *InternshipController) ApplicationsByStatus(c echo.Context) error {
	if err := ic.PastikanPenggunaAdmin(c); err != nil {
		return err
	}

	// Dapatkan status yang akan difilter dari permintaan (misalnya, "melamar", "ditolak", "diterima")
	status := c.QueryParam("status")

	// Query basis data untuk mendapatkan daftar aplikasi dengan status tertentu
	var applications []models.ApplicationStatus
	if err := ic.DB.Where("status = ?", status).Find(&applications).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mendapatkan aplikasi berdasarkan status"})
	}

	// Dapatkan daftar pengguna yang sesuai dengan aplikasi
	var users []models.User
	userIDs := []int{}
	for _, app := range applications {
		userIDs = append(userIDs, app.UserID)
	}

	if err := ic.DB.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mendapatkan pengguna berdasarkan aplikasi"})
	}

	// Return respons dengan daftar pengguna yang sesuai
	return c.JSON(http.StatusOK, users)
}


// CreateInternship digunakan untuk admin membuat lowongan magang.
func (ic *InternshipController) CreateInternship(c echo.Context) error {
	if err := ic.PastikanPenggunaAdmin(c); err != nil {
		return err
	}

	internship := new(models.Internship)
	if err := c.Bind(internship); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data tidak valid"})
	}

	if err := ic.DB.Create(internship).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membuat lowongan magang"})
	}

	return c.JSON(http.StatusCreated, internship)
}

// UpdateInternship digunakan untuk admin mengubah lowongan magang.
func (ic *InternshipController) UpdateInternship(c echo.Context) error {
	if err := ic.PastikanPenggunaAdmin(c); err != nil {
		return err
	}

	internshipID := c.Param("id")
	var internship models.Internship
	if err := ic.DB.Where("id = ?", internshipID).First(&internship).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Lowongan magang tidak ditemukan")
	}

	updatedInternship := new(models.Internship)
	if err := c.Bind(updatedInternship); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data pembaruan tidak valid"})
	}

	// Mengganti kolom yang diperlukan
	internship.Title = updatedInternship.Title
	internship.Description = updatedInternship.Description
	internship.Quota = updatedInternship.Quota

	if err := ic.DB.Save(&internship).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui lowongan magang"})
	}

	return c.JSON(http.StatusOK, internship)
}

// DeleteInternship digunakan untuk admin menghapus lowongan magang.
func (ic *InternshipController) DeleteInternship(c echo.Context) error {
	if err := ic.PastikanPenggunaAdmin(c); err != nil {
		return err
	}

	internshipID := c.Param("id")
	var internship models.Internship
	if err := ic.DB.Where("id = ?", internshipID).First(&internship).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Lowongan magang tidak ditemukan")
	}

	if err := ic.DB.Delete(&internship).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus lowongan magang"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Lowongan magang telah dihapus"})
}

// UpdateApplicationStatus digunakan untuk admin mengubah status pendaftaran user.
func (ic *InternshipController) UpdateApplicationStatus(c echo.Context) error {
	if err := ic.PastikanPenggunaAdmin(c); err != nil {
		return err
	}

	applicationID := c.Param("id")
	var applicationStatus models.ApplicationStatus
	if err := ic.DB.Where("id = ?", applicationID).First(&applicationStatus).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Status aplikasi tidak ditemukan")
	}

	updatedStatus := new(models.ApplicationStatus)
	if err := c.Bind(updatedStatus); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data pembaruan tidak valid"})
	}

	// Cek apakah status telah berubah menjadi "Diterima"
	if applicationStatus.Status != "Diterima" && updatedStatus.Status == "Diterima" {
		// Kirim email ke pengguna
		to := "alamat_email_pengguna@example.com"
		subject := "Selamat! Anda Diterima"
		message := "Selamat! Anda telah diterima untuk magang kami."

		// Konfigurasi email
		email := gomail.NewMessage()
		email.SetHeader("From", "pengirim@example.com")
		email.SetHeader("To", to)
		email.SetHeader("Subject", subject)
		email.SetBody("text/html", message)

		// Kirim email
		dialer := gomail.NewDialer("smtp.example.com", 587, "email_username", "email_password")
		if err := dialer.DialAndSend(email); err != nil {
			// Handle kesalahan jika gagal mengirim email
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengirim email"})
		}
	}

	// Selanjutnya, lanjutkan dengan perubahan status seperti yang telah Anda lakukan sebelumnya
	applicationStatus.Status = updatedStatus.Status

	if err := ic.DB.Save(&applicationStatus).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui status aplikasi"})
	}

	return c.JSON(http.StatusOK, applicationStatus)
}
