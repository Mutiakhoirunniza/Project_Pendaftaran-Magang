package controllers

import (
	"fmt"
	"log"
	"miniproject/constants"
	"miniproject/entity"
	"miniproject/infra/config"
	"miniproject/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// akun pengguna

// Fungsi RegisterUser digunakan untuk mendaftarkan pengguna baru.
func RegisterUser(c echo.Context) error {
	user := entity.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid user data",
			"error":   err.Error(),
		})
	}

	// Cek apakah pengguna sudah terdaftar berdasarkan alamat email
	var existingUser entity.User
	err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": constants.ErrUserAlreadyExists,
		})
	}

	// Jika pengguna belum terdaftar, simpan data pendaftaran ke dalam basis data
	err = config.DB.Create(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": constants.ErrFailedToRegister,
			"error":   err.Error(),
		})
	}

	// Mengirim respons HTTP berhasil setelah pengguna berhasil didaftarkan
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success create new user",
		"user":    user,
	})
}

// Fungsi LoginUserController digunakan untuk mengautentikasi pengguna dan memberikan token akses jika berhasil.
func LoginUserController(c echo.Context) error {
	// Membuat instance pengguna dan mengikat data dari permintaan HTTP
	user := entity.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Fail to parse request body",
			"error":   err.Error(),
		})
	}

	// Mencari pengguna dalam basis data berdasarkan email dan kata sandi
	err := config.DB.Where("email = ? AND password = ?", user.Email, user.Password).First(&user).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": constants.ErrFailedToLogIn,
			"error":   err.Error(),
		})
	}

	// Menghasilkan token akses untuk pengguna
	username := "user"
	role := "user"
	token, err := middleware.CreateToken(username, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": constants.ErrTokenCreationFailed,
			"error":   err.Error(),
		})
	}
	c.Set("user", token)
	UserResponse := entity.UserResponse{
		ID:    user.ID,
		Name:  user.Username,
		Email: user.Email,
		Token: token,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success login",
		"user":    UserResponse,
	})
}

// GetAllUsers digunakan untuk mendapatkan semua data pengguna.
func GetAllUsers(c echo.Context) error {
	var users []entity.User
	if err := config.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve users"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success: get all users",
		"users":   users,
	})
}

// GetUserByID digunakan untuk mendapatkan data pengguna berdasarkan ID.
func GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	var user entity.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success: get user by ID",
		"user":    user,
	})
}

// Fungsi UpdateUserByID digunakan untuk memperbarui data pengguna berdasarkan ID.
func UpdateUserByID(c echo.Context) error {
	// Mendapatkan ID pengguna dari parameter rute
	IdStr := c.Param("id")
	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		// Mengirim respons HTTP jika ID tidak valid
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	// Membuat instance baru dari entitas pengguna dan mengikat data dari permintaan HTTP
	user := new(entity.User)
	if err := c.Bind(user); err != nil {
		// Mengirim respons HTTP jika terjadi kesalahan dalam mengikat data pengguna
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Mencari pengguna yang ada dalam basis data berdasarkan ID
	var existingUser entity.User
	if err := config.DB.First(&existingUser, Id).Error; err != nil {
		// Mengirim respons HTTP jika pengguna tidak ditemukan
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Memperbarui data pengguna yang ada dengan data baru
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.Password = user.Password
	existingUser.Gender = user.Gender
	existingUser.PhoneNumber = user.PhoneNumber
	existingUser.UniversityName = user.UniversityName
	existingUser.UniversityAddress = user.UniversityAddress
	existingUser.Major = user.Major

	// Menyimpan perubahan data pengguna ke dalam basis data
	if err := config.DB.Save(&existingUser).Error; err != nil {
		// Mengirim respons HTTP jika terjadi kesalahan saat menyimpan perubahan
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Mengirim respons HTTP berhasil setelah pengguna diperbarui
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success update user",
		"user":    existingUser,
	})
}

// Menghapus data user berdasarkan ID
func DeleteUser(c echo.Context) error {
	IdStr := c.Param("id")
	Id, err := strconv.Atoi(IdStr)
	// Mengirim respons HTTP jika ID pengguna tidak valid
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid User ID")
	}

	var User entity.User
	// Mengirim respons HTTP jika pengguna tidak ditemukan
	if err := config.DB.First(&User, Id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User Not Found")
	}
	// Mengirim respons HTTP jika terjadi kesalahan saat menghapus pengguna
	if err := config.DB.Delete(&User).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success delete user",
		"User":    User,
	})
}

//  internship user

// GetInternshipListings digunakan untuk mendapatkan daftar lowongan magang
func GetInternshipListings(c echo.Context) error {
	var internshipListings []entity.Internship_Listing

	// Dapatkan semua daftar lowongan magang dari basis data
	if err := config.DB.Find(&internshipListings).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal mengambil daftar magang",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Daftar magang Terbaru", //diganti
		"listings": internshipListings,
	})
}

// ApplyForInternship ini digunakan untuk mengirimkan aplikasi pendaftaran magang
func ApplyForInternship(c echo.Context) error {
	// Deklarasi dan pengisian instansi ApplicationForm
	var formData entity.Internship_ApplicationForm
	fmt.Println("formData", formData)
	if err := c.Bind(&formData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Gagal mengikuti pendaftaran magang",
			"error":   err.Error(),
		})
	}

	// Mencari ID penawaran magang berdasarkan judul yang dipilih
	selectedTitle := formData.SelectedTitle
	var selectedListingID uint
	// sesuaikan cara mengambil data dari database
	var listing entity.Internship_Listing
	if err := config.DB.Where("title = ?", selectedTitle).First(&listing).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Penawaran magang tidak ditemukan",
		})
	}
	selectedListingID = listing.ID

	// Validasi data
	invalidData := make(map[string]string)
	if formData.Nim == "" {
		invalidData["nim"] = "Nim is required"
	}
	if formData.GPA <= 0 {
		invalidData["gpa"] = "GPA must be greater than 0"
	}
	if formData.EducationLevel == "" {
		invalidData["education_level"] = "Education level is required"
	}
	if formData.Username == "" {
		invalidData["username"] = "Username is required"
	}
	if formData.UserEmail == "" {
		invalidData["user_email"] = "User email is required"
	}

	if len(invalidData) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":     "Data formulir tidak valid",
			"invalidData": invalidData,
		})
	}

	// Simpan aplikasi ke dalam database
	application := entity.Internship_ApplicationForm{
		Model:               gorm.Model{},
		CV:                  formData.CV,
		Nim:                 formData.Nim,
		GPA:                 formData.GPA,
		EducationLevel:      formData.EducationLevel,
		UserID:              formData.UserID,
		Status:              "",
		UserEmail:           formData.UserEmail,
		Username:            formData.Username,
		SelectedTitle:       selectedTitle,
		InternshipListingID: selectedListingID,
		// Selected_Candidates: formData.Selected_Candidates,
	}
	log.Println("application", application)

	// Dapatkan penawaran magang berdasarkan selectedListingID
	var internshipListing entity.Internship_Listing
	if err := config.DB.First(&internshipListing, selectedListingID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Penawaran magang tidak ditemukan",
		})
	}

	if err := config.DB.Create(&application).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal memproses pendaftaran magang",
			"error":   err.Error(),
		})
	}

	// Kurangi kuota penawaran magang
	if internshipListing.Quota <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Kuota pendaftaran magang sudah penuh",
		})
	}
	internshipListing.Quota--
	if err := config.DB.Save(&internshipListing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal memproses pendaftaran magang",
			"error":   err.Error(),
		})
	}

	// Membuat dan menyimpan Selected_Candidate
	selectedCandidate := entity.Selected_Candidate{
		Model: gorm.Model{
			ID:        selectedListingID,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		InternshipApplicationFormID: application.ID,
		InternshipApplicationForm:   application,
	}

	if err := config.DB.Create(&selectedCandidate).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal memproses pendaftaran magang",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Pendaftaran magang berhasil disimpan",
	})
}

// CancelApplication digunakan untuk membatalkan formulir aplikasi berdasarkan ID.
func CancelApplication(c echo.Context) error {
	// Mendapatkan ID formulir aplikasi yang ingin dibatalkan 
	idParam := c.Param("id")

	// Mengonversi ID menjadi tipe data uint
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID tidak valid",
		})
	}

	// Cari formulir aplikasi berdasarkan ID
	var application entity.Internship_ApplicationForm
	if err := config.DB.Where("id = ?", id).First(&application).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Formulir aplikasi tidak ditemukan",
		})
	}

	// Periksa apakah formulir aplikasi sudah dibatalkan sebelumnya
	if application.IsCanceled {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Formulir aplikasi sudah dibatalkan sebelumnya",
		})
	}

	// Mengubah status formulir menjadi "Dibatalkan" dan menyimpan perubahan ke database
	application.Status = "Dibatalkan"
	application.IsCanceled = true
	if err := config.DB.Save(&application).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal membatalkan formulir aplikasi",
			"error":   err.Error(),
		})
	}

	// Mengembalikan kuota penawaran magang
	var internshipListing entity.Internship_Listing
	if err := config.DB.First(&internshipListing, application.InternshipListingID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Penawaran magang tidak ditemukan",
		})
	}

	internshipListing.Quota++
	if err := config.DB.Save(&internshipListing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal memproses pembatalan formulir aplikasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Formulir aplikasi berhasil dibatalkan",
	})
}

// GetApplicationStatus digunakan untuk mendapatkan status formulir aplikasi berdasarkan ID.
func GetApplicationStatus(c echo.Context) error {
	// Mendapatkan ID dari parameter URL
	idParam := c.Param("id")

	// Mengonversi ID menjadi tipe data uint
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID tidak valid",
		})
	}

	// Cari form aplikasi berdasarkan ID
	var application entity.Internship_ApplicationForm
	if err := config.DB.Where("id = ?", id).First(&application).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Form aplikasi tidak ditemukan",
		})
	}

	// Anda dapat mengakses status aplikasi melalui application.Status
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Status form aplikasi",
		"status":  application.Status,
	})
}
