package controllers

import (
	"context"
	"io"
	"log"
	"miniproject/entity"
	"miniproject/middleware"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB  *gorm.DB
	Log *log.Logger
}

func NewUserController(db *gorm.DB, logger *log.Logger) *UserController {
	return &UserController{DB: db, Log: logger}
}

// Register digunakan untuk mendaftarkan pengguna baru.
func (u *UserController) Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" || password == "" {
		u.Log.Println("Username and password are required") // Log pesan kesalahan
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username and password are required"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Printf("Error generating password hash: %v", err) // Log pesan kesalahan
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An error occurred"})
	}

	user := entity.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := u.DB.Create(&user).Error; err != nil {
		u.Log.Printf("Error creating user: %v", err) // Log pesan kesalahan
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An error occurred"})
	}
	u.Log.Println("User registered successfully") // Log pesan sukses
	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// Login digunakan untuk melakukan proses login pengguna.
func (u *UserController) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var user entity.User
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		u.Log.Printf("Login failed: %v", err) // Log pesan kesalahan
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Username or password is incorrect"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		u.Log.Printf("Login failed: %v", err) // Log pesan kesalahan
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Username or password is incorrect"})
	}
	u.Log.Println("Login successful") // Log pesan sukses
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}

// UpdateProfile digunakan untuk mengupdate profil pengguna.
func (u *UserController) UpdateProfile(c echo.Context) error {
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

	if email == "" || gender == "" || len(email) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email and gender are required, and email length is too large"})
	}

	if major == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Major is required"})
	}

	// Update the user's profile data in the database
	if err := u.DB.Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email":              email,
		"gender":             gender,
		"phone_number":       phoneNumber,
		"university_name":    universityName,
		"university_address": universityAddress,
		"major":              major,
	}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An error occurred"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

// UploadProfilePicture digunakan untuk mengunggah foto profil pengguna.
func (u *UserController) UploadProfilePicture(c echo.Context) error {
	file, err := c.FormFile("profile_picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Gagal mendapatkan file gambar")
	}
	// Periksa jenis file
	if !isImageFile(file.Filename) {
		return c.JSON(http.StatusBadRequest, "File yang diunggah harus berupa gambar")
	}
	// Periksa ukuran file (maksimal 3MB)
	maxSize := int64(3 * 1024 * 1024)
	if file.Size > maxSize {
		return c.JSON(http.StatusBadRequest, "Ukuran file terlalu besar (maksimal 3MB)")
	}

	// Baca file gambar
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Gagal membuka file gambar")
	}
	defer src.Close()

	// Inisialisasi koneksi ke GCS
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal terhubung ke GCS"})
	}
	defer client.Close()

	// Simpan file gambar di GCS
	bucketName := "nama-bucket-gcs" // Ganti dengan nama bucket GCS Anda
	objectName := "foto/" + file.Filename

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, src); err != nil {
		wc.Close()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengunggah foto ke GCS"})
	}
	wc.Close()

	// Simpan URL GCS ke database atau sesuai kebutuhan aplikasi Anda
	fotoURL := "https://storage.googleapis.com/" + bucketName + "/" + objectName

	return c.JSON(http.StatusOK, map[string]string{"image_url": fotoURL})
}

// Fungsi untuk memeriksa apakah file adalah gambar
func isImageFile(filename string) bool {
	// dalam format yang diizinkan.
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	extension := filepath.Ext(filename)
	for _, ext := range allowedExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}

// GetAllUsers digunakan untuk mendapatkan semua data pengguna.
func (u *UserController) GetAllUsers(c echo.Context) error {
	var users []entity.User
	if err := u.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve users"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success: get all users",
		"users":   users,
	})
}

// GetUserByID digunakan untuk mendapatkan data pengguna berdasarkan ID.
func (u *UserController) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	var user entity.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success: get user by ID",
		"user":    user,
	})
}

// DeleteUserByID digunakan untuk menghapus pengguna berdasarkan ID.
func (u *UserController) DeleteUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	var user entity.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	if err := u.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Success: user deleted"})
}

// GetInternshipListings digunakan untuk mendapatkan daftar lowongan magang.
func (u *UserController) GetInternshipListings(c echo.Context) error {
	var listings []entity.InternshipListing
	if err := u.DB.Find(&listings).Error; err != nil {
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
func (u *UserController) CheckIfUserFilledForm(c echo.Context, userID int) bool {
	var formData entity.InternshipApplicationForm
	if err := u.DB.Where("user_id = ?", userID).First(&formData).Error; err != nil {
		return false
	}
	return true
}

// HasPermission adalah fungsi bantu untuk memeriksa izin pengguna.
func HasPermission(userRole string) bool {
	if userRole == "User" {
		return true
	}
	return false
}

func (u *UserController) ChooseInternshipListing(c echo.Context) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token JWT tidak valid"})
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID pengguna tidak valid"})
	}

	if !u.CheckIfUserFilledForm(c, userIDInt) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Anda harus mengisi formulir pendaftaran terlebih dahulu"})
	}

	listingIDStr := c.Param("listing_id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID Pendaftaran Magang tidak valid"})
	}

	var listing entity.InternshipListing
	if err := u.DB.First(&listing, listingID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Pendaftaran Magang tidak ditemukan"})
	}

	if listing.Quota <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Kuota untuk pendaftaran magang ini sudah penuh"})
	}

	// Periksa izin pengguna (Anda perlu mendapatkan peran pengguna dari suatu sumber sebelumnya)
	userRole := "User"
	if HasPermission(userRole) {
		if listing.Quota > 0 {
			listing.SelectedCandidates = append(listing.SelectedCandidates, userIDInt)
			listing.Quota--
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Kuota untuk pendaftaran magang ini sudah penuh"})
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Anda tidak memiliki izin untuk memilih pendaftaran magang"})
	}

	if err := u.DB.Save(&listing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Terjadi kesalahan saat mengurangi kuota"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Pemilihan pendaftaran magang berhasil"})
}

// GetApplicationStatus digunakan untuk mendapatkan status aplikasi pengguna.
func (u *UserController) GetApplicationStatus(c echo.Context) error {
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
	if err := u.DB.Model(&entity.InternshipListing{}).Where("selected_candidates LIKE ?", "%"+strconv.Itoa(userIDInt)+"%").Find(&selectedListings).Error; err != nil {
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
func (u *UserController) CancelApplication(c echo.Context) error {
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
	if err := u.DB.First(&selectedListing, listingID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Listing yang dipilih tidak ditemukan"})
	}

	// Memeriksa apakah ID pengguna terdapat dalam daftar kandidat yang dipilih
	if !contains(selectedListing.SelectedCandidates, userIDInt) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Pengguna tidak terpilih untuk listing ini"})
	}

	// Memeriksa apakah status aplikasi sudah 'dibatalkan'
	if selectedListing.StatusPendaftaran == "dibatalkan" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Aplikasi sudah dibatalkan"})
	}

	// Mengubah status aplikasi menjadi 'dibatalkan' (Menunggu verifikasi admin)
	selectedListing.StatusPendaftaran = "dibatalkan"

	// Menyimpan perubahan pada listing ke dalam database
	if err := u.DB.Save(&selectedListing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membatalkan aplikasi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Permintaan pembatalan aplikasi berhasil. Menunggu verifikasi admin"})
}
