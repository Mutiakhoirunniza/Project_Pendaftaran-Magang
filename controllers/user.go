package controllers

import (
	"net/http"
	"project/models"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// UserController adalah tipe controller yang menangani operasi pengguna
type UserController struct {
    DB *gorm.DB
}

// NewUserController adalah fungsi untuk membuat instance UserController
func NewUserController(db *gorm.DB) *UserController {
    return &UserController{DB: db}
}

// RegisterUser digunakan untuk mendaftarkan pengguna.
func (uc *UserController) RegisterUser(c echo.Context) error {
    // Menerima data pendaftaran pengguna dari permintaan
    userData := new(models.User)
    if err := c.Bind(userData); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data tidak valid"})
    }

    // Validasi data pengguna
    if userData.Username == "" || userData.Password == "" || userData.Email == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username, Password, dan Email harus diisi"})
    }

    // Simpan pengguna ke database
    if err := uc.DB.Create(userData).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal menyimpan pengguna"})
    }

    // Return respons yang sesuai (contoh respons JSON)
    return c.JSON(http.StatusCreated, map[string]string{"message": "Pendaftaran berhasil"})
}

// LoginUser digunakan untuk mengotentikasi pengguna dan memberikan token JWT.
func (uc *UserController) LoginUser(c echo.Context) error {
    // Menerima data login pengguna dari permintaan
    username := c.FormValue("username")
    password := c.FormValue("password")

    // Cari pengguna berdasarkan username
    var user models.User
    if err := uc.DB.Where("username = ?", username).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Login gagal"})
    }

    // Verifikasi kata sandi
    if !checkPassword(password, user.Password) {
        return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Login gagal"})
    }

    // Generate JWT token (contoh)
    token, err := generateJWT(user.ID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Terjadi kesalahan saat menghasilkan token JWT"})
    }

    // Setelah berhasil terotentikasi, Anda dapat memberikan token JWT sebagai respons
    return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// GetUserProfile digunakan untuk mendapatkan profil pengguna yang sedang masuk.
func (uc *UserController) GetUserProfile(c echo.Context) error {
    // Dapatkan ID pengguna dari token JWT (contoh: token digambarkan sebagai string) atau sesi
    userID, ok := c.Get("userID").(string)
    if !ok {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
    }

    // Gunakan ID pengguna untuk mencari profil pengguna dari database
    var user models.User
    if err := uc.DB.Where("id = ?", userID).First(&user).Error; err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Profil pengguna tidak ditemukan")
    }

    // Return profil pengguna sebagai respons
    return c.JSON(http.StatusOK, user)
}

// GetUserProfileByID digunakan untuk mendapatkan profil pengguna berdasarkan ID pengguna tertentu.
func (uc *UserController) GetUserProfileByID(c echo.Context) error {
    // Dapatkan ID pengguna dari parameter URL
    userID := c.Param("id")

    // Gunakan ID pengguna untuk mencari profil pengguna dari database
    var user models.User
    if err := uc.DB.Where("id = ?", userID).First(&user).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Profil pengguna tidak ditemukan"})
    }

    // Return profil pengguna berdasarkan ID sebagai respons
    return c.JSON(http.StatusOK, user)
}

// UpdateUserProfile digunakan untuk mengganti profil pengguna.
func (uc *UserController) UpdateUserProfile(c echo.Context) error {
    // Dapatkan ID pengguna dari token JWT (contoh: token digambarkan sebagai string) atau sesi
    userID, ok := c.Get("userID").(string)
    if !ok {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
    }

    // Dapatkan data profil yang diperbarui dari permintaan
    updatedUser := new(models.User)
    if err := c.Bind(updatedUser); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data profil tidak valid"})
    }

    // Periksa apakah pengguna yang sedang mencoba mengganti profil adalah pemilik profil
    if userID != updatedUser.ID {
        return c.JSON(http.StatusForbidden, map[string]string{"message": "Anda tidak diizinkan mengganti profil pengguna lain"})
    }

    // Periksa apakah pengguna dengan ID yang sesuai ada di database
    var existingUser models.User
    if err := uc.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Profil pengguna tidak ditemukan"})
    }

    // Lakukan validasi data profil yang diperbarui di sini sesuai kebutuhan
    if updatedUser.Username == "" || updatedUser.Fullname == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Username dan Fullname harus diisi"})
    }

    // Perbarui profil pengguna di database
    existingUser.Username = updatedUser.Username
    existingUser.Fullname = updatedUser.Fullname
    // Update atribut lain yang diperlukan

    if err := uc.DB.Save(&existingUser).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Terjadi kesalahan saat menyimpan profil pengguna"})
    }

    // Return respons sukses
    return c.JSON(http.StatusOK, existingUser)
}

// DeleteUser digunakan untuk menghapus akun pengguna.
func (uc *UserController) DeleteUser(c echo.Context) error {
    // Dapatkan ID pengguna dari token JWT (contoh: token digambarkan sebagai string) atau sesi
    userID, ok := c.Get("userID").(string)
    if !ok {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
    }

    // Periksa apakah pengguna yang sedang mencoba menghapus akun adalah pemilik akun
    userToDeleteID := c.Param("id") // ID akun yang akan dihapus
    if userID != userToDeleteID {
        return c.JSON(http.StatusForbidden, map[string]string{"message": "Anda tidak diizinkan menghapus akun pengguna lain"})
    }

    // Periksa apakah pengguna dengan ID yang sesuai ada di database
    var existingUser models.User
    if err := uc.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Akun pengguna tidak ditemukan"})
    }

    // Hapus akun pengguna dari database
    if err := uc.DB.Delete(&existingUser).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Terjadi kesalahan saat menghapus akun pengguna"})
    }

    // Return respons sukses
    return c.JSON(http.StatusOK, map[string]string{"message": "Akun pengguna telah dihapus"})
}

// GetApplicationStatus digunakan untuk mendapatkan status aplikasi pengguna.
func (uc *UserController) GetApplicationStatus(c echo.Context) error {
    // Dapatkan ID pengguna dari token JWT atau sesi
    userID, ok := c.Get("userID").(string)
    if !ok {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
    }

    // Dapatkan daftar status aplikasi pengguna dari database
    var statusList []models.ApplicationStatus
    if err := uc.DB.Where("user_id = ?", userID).Find(&statusList).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil status aplikasi"})
    }

    // Buat daftar respons yang menggambarkan status setiap aplikasi
    var responseList []map[string]interface{}
    for _, status := range statusList {
        applicationStatus := map[string]interface{}{
            "internship_id": status.InternshipID,
            "status":        status.Status,
        }
        responseList = append(responseList, applicationStatus)
    }

    // Return daftar status aplikasi sebagai respons
    return c.JSON(http.StatusOK, responseList)
}

// CancelApplication digunakan untuk pembatalan aplikasi jika status masih dalam tahap pengajuan.
func (uc *UserController) CancelApplication(c echo.Context) error {
    // Dapatkan ID pengguna dari token JWT atau sesi
    userID, ok := c.Get("userID").(string)
    if !ok {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
    }

    // Dapatkan ID aplikasi dari permintaan
    applicationID := c.Param("id")

    // Cek status aplikasi
    var applicationStatus models.ApplicationStatus
    if err := uc.DB.Where("id = ? AND user_id = ?", applicationID, userID).First(&applicationStatus).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Aplikasi tidak ditemukan"})
    }

    // Hanya izinkan pembatalan jika status adalah "Pengajuan"
    if applicationStatus.Status != "Pengajuan" {
        return c.JSON(http.StatusForbidden, map[string]string{"message": "Anda tidak dapat membatalkan aplikasi ini"})
    }

    // Hapus aplikasi dari database
    if err := uc.DB.Delete(&applicationStatus).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal membatalkan aplikasi"})
    }

    // Return respons yang sesuai
    return c.JSON(http.StatusOK, map[string]string{"message": "Pembatalan aplikasi berhasil"})
}

