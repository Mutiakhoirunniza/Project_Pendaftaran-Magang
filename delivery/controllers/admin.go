package controllers

import (
	"miniproject/constants"
	"miniproject/entity"
	"miniproject/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

// NewAdminController membuat instance AdminController
func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

// proses login admin
func (a *AdminController) LoginAdmin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if a.authenticateAdmin(username, password) {
		// Login berhasil, menghasilkan token JWT
		token, err := middleware.GenerateJWTToken(username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Gagal menghasilkan token")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Login berhasil",
			"token":   token,
		})
	}
	return c.JSON(http.StatusUnauthorized, "Login gagal, Periksa kembali username dan password Anda.")
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

// Membuat formulir pendaftaran lowongan magang
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

    // Simpan formulir pendaftaran ke database
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
