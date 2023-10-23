package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"miniproject/constants"
	"miniproject/entity"
	"miniproject/middleware"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type AdminController struct {
	DB     *gorm.DB
	Logger *log.Logger
}

func NewAdminController(db *gorm.DB, logger *log.Logger) *AdminController {
	return &AdminController{DB: db, Logger: logger}
}

func (a *AdminController) LoginAdmin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if a.authenticateAdmin(username, password) {
		token, err := middleware.GenerateJWTToken(username)
		if err != nil {
			a.Logger.Printf("Gagal menghasilkan token: %s", err)
			return c.JSON(http.StatusInternalServerError, "Gagal menghasilkan token")
		}
		a.Logger.Printf("Login berhasil: %s", username)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Login berhasil",
			"token":   token,
		})
	} else {
		a.Logger.Printf("Login gagal: %s", username)
		a.Logger.Print("Login gagal, Periksa kembali username dan password Anda.")
		return c.JSON(http.StatusUnauthorized, "Login gagal, Periksa kembali username dan password Anda.")
	}
}


// authenticateAdmin memeriksa apakah username dan password admin sesuai
func (a *AdminController) authenticateAdmin(username, password string) bool {
	return username == "admin" && password == "password"
}

// Mengambil data admin berdasarkan ID
func (a *AdminController) GetAdminByID(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	var admin entity.Admin
	if err := a.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak ditemukan")
	}

	return c.JSON(http.StatusOK, admin)
}

// Mengubah data admin berdasarkan ID
func (a *AdminController) UpdateAdmin(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	var admin entity.Admin
	if err := a.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak dapat ditemukan")
	}

	if err := a.DB.Save(&admin).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Gagal menyimpan perubahan")
	}

	return c.JSON(http.StatusOK, admin)
}

// Menghapus data admin berdasarkan ID
func (a *AdminController) DeleteAdmin(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	var admin entity.Admin
	if err := a.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak ditemukan")
	}

	if err := a.DB.Delete(&admin).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Gagal menghapus Admin")
	}

	return c.JSON(http.StatusOK, "Admin berhasil dihapus")
}

// Membuat daftar lowongan magang baru
func (a *AdminController) CreateInternshipListing(c echo.Context) error {
	var inputData entity.InputData
	if err := c.Bind(&inputData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data tidak valid"})
	}

	// Validasi data
	if inputData.Quota <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Kuota harus lebih dari 0"})
	}

	// Inisialisasi CreatedDate dengan waktu saat ini
	listing := entity.InternshipListing{
		Title:       inputData.Title,
		Description: inputData.Description,
		Quota:       inputData.Quota,
		CreatedDate: time.Now(),
	}

	// Simpan daftar lowongan magang ke database
	if err := a.DB.Create(&listing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membuat daftar lowongan magang"})
	}

	return c.JSON(http.StatusCreated, listing)
}

func (a *AdminController) CreateInternshipApplicationForm(c echo.Context) error {
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
	if file.Size > 3*1024*1024 { // 3 MB dalam byte
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
	bucketName := "nama-bucket-gcs" // Ganti dengan nama bucket GCS Anda
	objectName := "cv/" + file.Filename

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, cv); err != nil {
		wc.Close()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengunggah CV ke GCS"})
	}
	if err := wc.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal menutup writer GCS"})
	}

	// Setelah berhasil menyimpan di GCS, Anda dapat menyimpan URL GCS di database atau yang sesuai.

	// Simpan formulir pendaftaran ke database (termasuk URL GCS CV)
	formData.CVURL = "https://storage.googleapis.com/" + bucketName + "/" + objectName

	if err := a.DB.Create(&formData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membuat formulir pendaftaran"})
	}

	return c.JSON(http.StatusCreated, formData)
}

// Mengubah status pendaftaran berdasarkan ID formulir pendaftaran
func (a *AdminController) UpdateApplicationStatus(c echo.Context) error {
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
	if err := a.DB.Where("internship_application_form_id = ?", formID).First(&existingStatus).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Status pendaftaran tidak ditemukan")
	}

	// Update status
	existingStatus.Status = statusData.Status

	if err := a.DB.Save(&existingStatus).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengubah status pendaftaran"})
	}

	return c.JSON(http.StatusOK, existingStatus)
}

// Fungsi untuk memverifikasi atau mengizinkan otomatis status pembatalan pengguna
func (a *AdminController) VerifyCancelApplication(c echo.Context) error {
	// Ambil ID formulir pendaftaran dari parameter URL
	formID, err := strconv.Atoi(c.Param("formID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Formulir Pendaftaran tidak valid")
	}

	// Cari formulir pendaftaran berdasarkan ID
	var formData entity.InternshipApplicationForm
	if err := a.DB.First(&formData, formID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Formulir pendaftaran tidak ditemukan")
	}

	// Cek apakah status formulir pendaftaran saat ini adalah "Pengajuan"
	if formData.Status != constants.StatusPending {
		return c.JSON(http.StatusBadRequest, "Hanya formulir dalam status 'Pengajuan' yang dapat diverifikasi atau diizinkan otomatis")
	}

	// Ubah status formulir pendaftaran menjadi "Dibatalkan" secara otomatis
	formData.Status = constants.StatusCanceled

	// Simpan perubahan status formulir pendaftaran ke database
	if err := a.DB.Save(&formData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Gagal mengubah status pendaftaran")
	}

	return c.JSON(http.StatusOK, formData)
}

func SendEmailToUser(userEmail, name, status string) error {
	if status != "diterima" {
		return nil // Jika status bukan "diterima", tidak perlu mengirim email
	}
	// Mendapatkan tanggal saat ini
	currentDate := time.Now()
	// Menambahkan 1 hari ke tanggal saat ini
	startDate := currentDate.AddDate(0, 0, 1).Format("02/01/2006")
	// Menambahkan 3 bulan ke "Mulai Tanggal" 
	endDate := currentDate.AddDate(0, 3, 0).Format("02/01/2006")

	smtpServer := os.Getenv("SMTPSERVER")
	smtpPortStr := os.Getenv("SMTPPORT")
	smtpUsername := os.Getenv("SMTPUSERNAME")
	smtpPassword := os.Getenv("SMTPPASSWORD")

	sender := smtpUsername
	recipient := userEmail
	subject := "Selamat Bergabung Sebagai Magang"
	emailBody := `
    <!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Selamat Bergabung Sebagai Magang di PT. Krisnadwipayana</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
            color: #333;
        }

        p {
            margin: 10px 0;
        }

        ul, ol {
            margin: 10px 0;
            padding-left: 20px;
        }

        strong {
            font-weight: bold;
        }

        a {
            color: #007bff;
            text-decoration: none;
        }

        a:hover {
            text-decoration: underline;
        }

        .container {
            background-color: #fff;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }

        .header {
            background-color: #007bff;
            color: #fff;
            text-align: center;
            padding: 10px;
        }

        .footer {
            background-color: #007bff;
            color: #fff;
            text-align: center;
            padding: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Selamat Bergabung Sebagai Magang</h1>
        </div>
        <p>Halo, ` + name + `</p>
        
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

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	msg := []byte("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n" +
		"\r\n" +
		emailBody)

	err := smtp.SendMail(smtpServer+":"+smtpPortStr, auth, sender, []string{recipient}, msg)
	if err != nil {
		return fmt.Errorf("Gagal mengirim email: %w", err)
	}

	return nil
}
