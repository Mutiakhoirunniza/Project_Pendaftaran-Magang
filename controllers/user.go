package controllers

import (
	"context"
	"io"
	"miniproject/entity"
	"miniproject/infra/config"
	"miniproject/middleware"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
)

// Register digunakan untuk mendaftarkan pengguna baru.
func Register(c echo.Context) error {
	user := entity.User{}
	c.Bind(&user)

	if err := config.DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new user",
		"user":    user,
	})
}

// Login digunakan untuk melakukan proses login pengguna.
func LoginUserController(c echo.Context) error {
	User := entity.User{}
	c.Bind(&User)

	err := config.DB.Where("email = ? AND password = ?", User.Email, User.Password).First(&User).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Fail login",
			"error":   err.Error(),
		})
	}
	token, err := middleware.GenerateJWTToken(User.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Fail login",
			"error":   err.Error(),
		})
	}

	UserResponse := entity.UserResponse{
		ID:      int(User.ID),
		Username: User.Username,
		Email:    User.Email,
		Token:    token,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success login",
		"user":    UserResponse,
	})
}

// UpdateProfile digunakan untuk mengupdate profil pengguna.
func UpdateProfileAndUploadPicture(c echo.Context) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid JWT token"})
	}

	email := c.FormValue("email")
	gender := c.FormValue("gender")
	phoneNumber := c.FormValue("phone_number")
	universityName := c.FormValue("university_name")
	universityAddress := c.FormValue("university_address")
	major := c.FormValue("major")

	// Validasi input
	if email == "" || gender == "" || major == "" || len(email) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input data"})
	}

	// Update the user's profile data in the database
	user := &entity.User{}
	if err := config.DB.Model(user).Where("id = ?", userID).Updates(map[string]interface{}{
		"email":              email,
		"gender":             gender,
		"phone_number":       phoneNumber,
		"university_name":    universityName,
		"university_address": universityAddress,
		"major":              major,
	}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An error occurred while updating profile"})
	}

	// Handle profile picture upload
	file, err := c.FormFile("profile_picture")
	if err == nil {
		// Periksa jenis file
		ext := strings.ToLower(filepath.Ext(file.Filename))
        if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
            return c.JSON(http.StatusBadRequest, "Uploaded file must be an image")
        }
		// Periksa ukuran file (maksimal 3MB)
		maxSize := int64(3 * 1024 * 1024)
		if file.Size > maxSize {
			return c.JSON(http.StatusBadRequest, "File size is too large (maximum 3MB)")
		}

		// Baca file gambar
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to open image file")
		}
		defer src.Close()

		// Inisialisasi koneksi ke GCS
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to connect to GCS"})
		}
		defer client.Close()

		// Simpan file gambar di GCS
		bucketName := "krisnadwipayana" // Ganti dengan nama bucket GCS Anda
		objectName := "foto/" + file.Filename

		wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
		if _, err := io.Copy(wc, src); err != nil {
			wc.Close()
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to upload photo to GCS"})
		}
		wc.Close()

		// Simpan URL GCS ke dalam entitas User
		foto := "https://storage.googleapis.com/" + bucketName + "/" + objectName
		user.ProfilePicture = foto

		// Simpan perubahan pada entitas User ke dalam database
		if err := config.DB.Save(user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save profile picture URL"})
		}

		// Return the URL of the uploaded profile picture
		return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully", "image_url": foto})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
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

// DeleteUserByID digunakan untuk menghapus pengguna berdasarkan ID.
func DeleteUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	var user entity.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Success: user deleted"})
}

// GetInternshipListings digunakan untuk mendapatkan daftar lowongan magang.
func GetInternshipListings(c echo.Context) error {
	var listings []entity.InternshipListing
	if err := config.DB.Find(&listings).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve internship listings"})
	}

	var result []map[string]interface{}
	for _, listing := range listings {
		listingData := map[string]interface{}{
			"id":          listing.ID,
			"title":       listing.Title,
			"description": listing.Description,
			"quota":       listing.Quota,
		}
		result = append(result, listingData)
	}

	return c.JSON(http.StatusOK, result)
}

// CheckIfUserFilledForm digunakan untuk memeriksa apakah pengguna sudah mengisi formulir pendaftaran magang.
func CheckIfUserFilledForm(_ echo.Context, userID int) bool {
	err := config.DB.Where("user_id = ?", userID).First(&entity.InternshipApplicationForm{}).Error
	if err != nil {
		// Jika data tidak ditemukan, kembalikan false
		return false
	}
	// Jika data ditemukan, kembalikan true
	return true
}


// HasPermission adalah fungsi bantu untuk memeriksa izin pengguna.
func HasPermission(userRole string) bool {
	if userRole == "User" {
		return true
	}
	return false
}

func ChooseInternshipListing(c echo.Context) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token JWT tidak valid"})
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID pengguna tidak valid"})
	}

	if !CheckIfUserFilledForm(c, userIDInt) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Anda harus mengisi formulir pendaftaran terlebih dahulu"})
	}

	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID Pendaftaran Magang tidak valid"})
	}

	var listing entity.InternshipListing
	if err := config.DB.First(&listing, listingID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Pendaftaran Magang tidak ditemukan"})
	}

	if listing.Quota <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Kuota untuk pendaftaran magang ini sudah penuh"})
	}

	// Periksa izin pengguna (Anda perlu mendapatkan peran pengguna dari suatu sumber sebelumnya)
	userRole := "User"
	if HasPermission(userRole) {
		if listing.Quota > 0 {
			listing.SelectedCandidates = append(listing.SelectedCandidates, entity.SelectedCandidate{CandidateID: userIDInt})
			listing.Quota--
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Kuota untuk pendaftaran magang ini sudah penuh"})
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Anda tidak memiliki izin untuk memilih pendaftaran magang"})
	}

	// Simpan perubahan dalam database
	if err := config.DB.Save(&listing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Terjadi kesalahan saat mengurangi kuota"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Pemilihan pendaftaran magang berhasil"})
}

// GetApplicationStatus digunakan untuk mendapatkan status aplikasi pengguna.
func GetApplicationStatus(c echo.Context) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid JWT token"})
	}

	// Parse user ID as an integer
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	// Query the database to get the internship listings chosen by the user
	var selectedListings []entity.InternshipListing
	if err := config.DB.Model(&entity.InternshipListing{}).Where("selected_candidates LIKE ?", "%"+strconv.Itoa(userIDInt)+"%").Find(&selectedListings).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve selected listings"})
	}

	// Prepare the response data
	var applicationStatus []map[string]interface{}
	for _, listing := range selectedListings {
		listingData := map[string]interface{}{
			"listing_id":  listing.ID,
			"title":       listing.Title,
			"description": listing.Description,
			"quota":       listing.Quota,
		}
		applicationStatus = append(applicationStatus, listingData)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "Success: get application status",
		"selected_listings": applicationStatus,
	})
}

// contains adalah fungsi bantu untuk memeriksa apakah elemen ada dalam slice.
func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// CancelApplication digunakan untuk membatalkan aplikasi pengguna.
func CancelApplication(c echo.Context) error {
	// Mendapatkan ID pengguna dari token JWT
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token JWT tidak valid"})
	}

	// Mengurai ID pengguna ke dalam bentuk integer
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID pengguna tidak valid"})
	}

	// Mengurai ID listing dari permintaan
	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID listing tidak valid"})
	}

	// Mengambil data listing yang dipilih dari database
	var selectedListing entity.InternshipListing
	if err := config.DB.Preload("SelectedCandidates").First(&selectedListing, listingID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Listing yang dipilih tidak ditemukan"})
	}

	// Membuat fungsi utilitas untuk memeriksa apakah userIDInt terdapat dalam selectedListing.SelectedCandidates
	containsUserID := func() bool {
		for _, candidate := range selectedListing.SelectedCandidates {
			if candidate.CandidateID == userIDInt {
				return true
			}
		}
		return false
	}

	// Memeriksa apakah ID pengguna terdapat dalam daftar kandidat yang dipilih
	if !containsUserID() {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Pengguna tidak terpilih untuk listing ini"})
	}

	// // Memeriksa apakah ID pengguna terdapat dalam daftar kandidat yang dipilih
	// if !contains(selectedListing.SelectedCandidates, userIDInt) {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"message": "Pengguna tidak terpilih untuk listing ini"})
	// }

	// Memeriksa apakah status aplikasi sudah 'dibatalkan'
	if selectedListing.StatusPendaftaran == "dibatalkan" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Aplikasi sudah dibatalkan"})
	}

	// Mengubah status aplikasi menjadi 'dibatalkan' (Menunggu verifikasi admin)
	selectedListing.StatusPendaftaran = "dibatalkan"

	// Menyimpan perubahan pada listing ke dalam database
	if err := config.DB.Save(&selectedListing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membatalkan aplikasi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Permintaan pembatalan aplikasi berhasil. Menunggu verifikasi admin"})
}
