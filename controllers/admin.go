package controllers

import (
	"context"
	"fmt"
	"io"
	"miniproject/constants"
	"miniproject/entity"
	"miniproject/infra/config"
	"miniproject/middleware"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

func LoginAdminController(c echo.Context) error {
	admin := entity.Admin{}
	c.Bind(&admin)

	err := config.DB.Where("email = ? AND password = ?", admin.Email, admin.Password).First(&admin).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Fail login",
			"error":   err.Error(),
		})
	}
	token, err := middleware.GenerateJWTToken(admin.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Fail login",
			"error":   err.Error(),
		})
	}

	AdminResponse := entity.AdminResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
		Token:    token,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login",
		"admin":   AdminResponse,
	})
}

// authenticateAdmin memeriksa apakah username dan password admin sesuai
func authenticateAdmin(username, password string) bool {
	return username == "admin" && password == "password"
}

// Mengambil data admin berdasarkan ID
func GetAdminByID(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	var admin entity.Admin
	if err := config.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak ditemukan")
	}

	return c.JSON(http.StatusOK, admin)
}

// Mengubah data admin berdasarkan ID
func UpdateAdmin(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	var admin entity.Admin
	if err := config.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak dapat ditemukan")
	}

	if err := config.DB.Save(&admin).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Gagal menyimpan perubahan")
	}

	return c.JSON(http.StatusOK, admin)
}

// Menghapus data admin berdasarkan ID
func DeleteAdmin(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	var admin entity.Admin
	if err := config.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak ditemukan")
	}

	if err := config.DB.Delete(&admin).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Gagal menghapus Admin")
	}

	return c.JSON(http.StatusOK, "Admin berhasil dihapus")
}

// Membuat daftar lowongan magang baru
func CreateInternshipListing(c echo.Context) error {
	var inputData entity.InternshipListing
	if err := c.Bind(&inputData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data tidak valid"})
	}

	// Validasi data
	if inputData.Quota <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Kuota harus lebih dari 0"})
	}

	// Inisialisasi CreatedDate dengan waktu saat ini
	inputData.CreatedDate = time.Now()

	// Simpan daftar lowongan magang ke database
	if err := config.DB.Create(&inputData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membuat daftar lowongan magang"})
	}

	return c.JSON(http.StatusCreated, inputData)
}

func CreateInternshipApplicationForm(c echo.Context) error {
	var formData entity.InternshipApplicationForm
	if err := c.Bind(&formData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data formulir tidak valid"})
	}

	// Validasi data formulir
	if formData.InternshipListingID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID Daftar Lowongan Magang harus lebih dari 0"})
	}

	// Validasi ukuran file CV (maksimal 3 MB)
	file, err := c.FormFile("cv")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "CV tidak dapat diunggah"})
	}
	if file.Size > 3*1024*1024 { // 3 MB
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "CV melebihi batas ukuran maksimal (3 MB)"})
	}

	// Dapatkan file CV dari form
	cv, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membuka file CV"})
	}
	defer cv.Close()

	// Inisialisasi koneksi ke GCS
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal terhubung ke GCS"})
	}
	defer client.Close()

	// Simpan file CV di GCS
	bucketName := "krisnadwipayana"
	objectName := "cv/" + file.Filename

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, cv); err != nil {
		wc.Close()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengunggah CV ke GCS"})
	}
	if err := wc.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal menutup writer GCS"})
	}

	// Simpan formulir pendaftaran ke database (termasuk URL GCS CV)
	formData.CV = "https://storage.googleapis.com/" + bucketName + "/" + objectName

	if err := config.DB.Create(&formData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membuat formulir pendaftaran"})
	}

	return c.JSON(http.StatusCreated, formData)
}

// Mengubah status pendaftaran berdasarkan ID formulir pendaftaran
func UpdateApplicationStatus(c echo.Context) error {
	formID, err := strconv.Atoi(c.Param("formID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Formulir Pendaftaran tidak valid")
	}

	var statusData entity.ApplicationStatus
	if err := c.Bind(&statusData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data status tidak valid"})
	}

	// Validasi data status
	validStatusValues := []string{
		constants.StatusPending,
		constants.StatusVerified,
		constants.StatusAccepted,
		constants.StatusRejected,
		constants.StatusCanceled,
	}
	isValidStatus := false
	for _, validStatus := range validStatusValues {
		if statusData.Status == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Status tidak valid"})
	}

	// Cari status berdasarkan ID formulir pendaftaran
	var existingStatus entity.ApplicationStatus
	if err := config.DB.Where("internship_application_form_id = ?", formID).First(&existingStatus).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Status pendaftaran tidak ditemukan")
	}

	// Update status
	existingStatus.Status = statusData.Status

	if err := config.DB.Save(&existingStatus).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengubah status pendaftaran"})
	}

	return c.JSON(http.StatusOK, existingStatus)
}

// Fungsi untuk memverifikasi atau mengizinkan otomatis status pembatalan pengguna
func VerifyCancelApplication(c echo.Context) error {
	formID, err := strconv.Atoi(c.Param("formID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Formulir Pendaftaran tidak valid")
	}

	var formData entity.InternshipApplicationForm
	if err := config.DB.First(&formData, formID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Formulir pendaftaran tidak ditemukan")
	}

	// Cek apakah status formulir pendaftaran saat ini adalah "Pengajuan"
	if formData.Status != constants.StatusPending {
		return c.JSON(http.StatusBadRequest, "Hanya formulir dalam status 'Pengajuan' yang dapat diverifikasi atau diizinkan otomatis")
	}

	// Admin dapat memutuskan untuk mengizinkan atau menolak pembatalan formulir
	approved, err := strconv.ParseBool(c.QueryParam("approved"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Parameter 'approved' tidak valid")
	}

	if approved {
		// Jika diizinkan, ubah status formulir pendaftaran menjadi "Dibatalkan"
		formData.Status = constants.StatusCanceled

		// Simpan perubahan status formulir pendaftaran ke database
		if err := config.DB.Save(&formData).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, "Gagal mengubah status pendaftaran")
		}
	}

	return c.JSON(http.StatusOK, formData)
}

func SendEmailToUser(userEmail, username, status string) error {
	// Hanya kirim email jika status adalah "diterima"
	if status != "diterima" {
		return nil
	}
	
	// Mengambil tanggal saat ini
	currentDate := time.Now()
	// Menambahkan 1 hari ke tanggal saat ini
	startDate := currentDate.AddDate(0, 0, 1).Format("02/01/2006")
	// Menambahkan 3 bulan ke "Mulai Tanggal"
	endDate := currentDate.AddDate(0, 3, 0).Format("02/01/2006")

	// Mengambil konfigurasi server SMTP dari variabel lingkungan
	smtpServer := os.Getenv("SMTPSERVER")
	smtpPortStr := os.Getenv("SMTPPORT")
	smtpUsername := os.Getenv("SMTPUSERNAME")
	smtpPassword := os.Getenv("SMTPPASSWORD")

	// Konversi smtpPortStr menjadi int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	// Pesan email dalam format HTML
	message := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Selamat Bergabung Sebagai Magang di PT. Krisnadwipayana</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f3f3f3;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					padding: 20px;
					background-color: #ffffff;
					box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #0073e6;
				}
				p {
					font-size: 16px;
				}
				ul {
					list-style-type: disc;
				}
				ol {
					list-style-type: decimal;
				}
				.footer {
					background-color: #0073e6;
					color: #fff;
					padding: 10px 0;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div "header">
					<h1>Selamat Bergabung Sebagai Magang</h1>
				</div>
				<p>Halo, ` + username + `</p>

				<p>Kami dengan senang hati ingin memberitahukan bahwa Anda telah diterima sebagai magang di PT. Krisnadwipayana. Kami sangat bersemangat untuk menyambut Anda menjadi bagian dari tim kami. Keputusan ini didasari oleh potensi dan kualifikasi yang luar biasa yang Anda tunjukkan selama proses seleksi.</p>

				<p><strong>Detail Kontrak Magang:</strong></p>
				<ul>
					<li>Mulai Tanggal: ` + startDate + `</li>
					<li>Berakhir Tanggal: ` + endDate + `</li>
					<li>Jam Kerja: 08.00 WIB</li>
					<li>Lokasi: JL. Gatot Utomo Kav 5</li>
				</ul>

				<p><strong>Langkah Selanjutnya:</strong></p>
				<ol>
					<li>Konfirmasi Kehadiran: Tolong konfirmasi penerimaan ini dengan membalas email ini atau menghubungi Mutiakhoirunniza di <a href="mailto:mutiakhoirunniza@ac.id">mutiakhoirunniza@ac.id</a> atau 08316281026.</li>
					<li>Dokumen-dokumen: Kami akan mengirimkan Anda semua dokumen dan formulir yang perlu Anda isi sebagai persyaratan magang. Mohon lengkapi dan kembalikan dokumen-dokumen ini sesegera mungkin.</li>
				</ol>

				<p><strong>Pengenalan Tim:</strong> Anda akan dikenalkan kepada tim yang akan menjadi mentormu selama magang. Mereka akan membantu Anda beradaptasi dan memberikan bimbingan selama Anda belajar di perusahaan.</p>

				<p>Kami sangat menghargai kesempatan ini untuk bekerja bersama Anda dan membantu Anda dalam pengembangan karier Anda. Jika Anda memiliki pertanyaan atau perlu bantuan lebih lanjut, jangan ragu untuk menghubungi kami kapan saja.</p>

				<p>Kami berharap Anda akan mendapatkan pengalaman yang berharga selama magang di PT. Krisnadwipayana dan berkembang bersama kami. Selamat datang di tim kami!</p>

				<p>Terima kasih dan selamat bersiap-siap untuk memulai petualangan baru Anda.</p>
				<div class="footer">
					<p>Salam, Mutia Khoirunniza</p>
					<p>JL. Gatot Utomo Kav 5</p>
					<p>08316281026</p>
				</div>
			</div>
		</body>
		</html>
	`

	// Konfigurasi pengiriman email menggunakan gomail
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUsername)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Selamat Bergabung Sebagai Magang")
	m.SetBody("text/html", message)

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("Gagal mengirim email: %w", err)
	}

	return nil
}