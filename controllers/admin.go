package controllers

import (
	"fmt"
	"miniproject/constants"
	"miniproject/entity"
	"miniproject/helpers"
	"miniproject/infra/config"
	"miniproject/middleware"
	"net/http"
	"strconv"
	

	"github.com/labstack/echo/v4"
)


// pengguna akun

// Fungsi LoginAdminController digunakan untuk mengautentikasi admin dan memberikan token akses jika berhasil.
func LoginAdminController(c echo.Context) error {
	admin := entity.Admin{}
	if err := c.Bind(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Fail to parse request body",
			"error":   err.Error(),
		})
	}

	// Mencari admin dalam basis data berdasarkan alamat email dan kata sandi
	err := config.DB.Where("email = ? AND password = ?", admin.Email, admin.Password).First(&admin).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Failed to log in",
			"error":   err.Error(),
		})
	}

	// Menghasilkan token akses untuk admin
	username := "admin"
	role := "admin"
	token, err := middleware.CreateToken(username, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
	}
	c.Set("admin", token) // Menyimpan token dalam konteks

	AdminResponse := entity.AdminResponse{
		ID:    admin.ID,
		Name:  admin.Username,
		Email: admin.Email,
		Token: token,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success login",
		"user":    AdminResponse,
	})
}
// Fungsi GetAdminByID digunakan untuk mengambil data admin berdasarkan ID.
func GetAdminByID(c echo.Context) error {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID Admin tidak valid")
	}

	// Mencari admin dalam basis data berdasarkan ID
	var admin entity.Admin
	if err := config.DB.First(&admin, ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Admin tidak ditemukan")
	}

	// Mengirim respons HTTP berhasil dengan data admin yang ditemukan
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"admin":   admin,
	})
}
// Fungsi UpdateAdminController digunakan untuk mengupdate data admin berdasarkan ID.
func UpdateAdminController(c echo.Context) error {
	IdStr := c.Param("id")
	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	// Membuat instance baru dari entitas admin dan mengikat data dari permintaan HTTP
	admin := new(entity.Admin)
	if err := c.Bind(admin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Mencari admin yang ada dalam basis data berdasarkan ID
	var existingAdmin entity.Admin
	if err := config.DB.First(&existingAdmin, Id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Admin not found")
	}

	// Memperbarui data admin yang ada dengan data baru dari permintaan
	existingAdmin.Username = admin.Username
	existingAdmin.Email = admin.Email
	existingAdmin.Password = admin.Password

	// Menyimpan perubahan data admin ke dalam basis data
	if err := config.DB.Save(&existingAdmin).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Mengirim respons HTTP berhasil setelah admin diperbarui
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success update admin",
		"admin":   existingAdmin,
	})
}


// semua data internships admin

// Membuat lowongan magang baru
func CreateInternshipListing(c echo.Context) error {
	// Bind data lowongan dari request body
	listing := entity.Internship_Listing{}
	if err := c.Bind(&listing); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Gagal mem-parsing request body",
			"error":   err.Error(),
		})
	}

	// Lakukan validasi kuota
	if listing.Quota <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Kuota harus lebih dari 0",
		})
	}

	if err := config.DB.Create(&listing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal membuat lowongan magang",
			"error":   err.Error(),
		})
	}

	// Cetak daftar lowongan magang setelah pembuatan
	var internshipListings []entity.Internship_Listing
	if err := config.DB.Find(&internshipListings).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal mengambil daftar magang",
			"error":   err.Error(),
		})
	}

	fmt.Printf("Daftar Magang : %+v\n", internshipListings)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Lowongan magang berhasil dibuat",
		"listing": listing,
	})
}
// Memperbarui lowongan magang berdasarkan ID
func UpdateInternshipListingByID(c echo.Context) error {
	id := c.Param("id")

	// Bind data lowongan dari request body
	listing := entity.Internship_Listing{}
	if err := c.Bind(&listing); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Gagal mem-parsing request body",
			"error":   err.Error(),
		})
	}

	if err := config.DB.Model(&entity.Internship_Listing{}).Where("id = ?", id).Updates(&listing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal memperbarui lowongan magang",
			"error":   err.Error(),
		})
	}

	// Cetak data lowongan magang setelah pembaruan
	var updatedListing entity.Internship_Listing
	if err := config.DB.Where("id = ?", id).First(&updatedListing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal mengambil data lowongan magang yang diperbarui",
			"error":   err.Error(),
		})
	}

	fmt.Printf("Data Lowongan yang Diperbarui: %+v\n", updatedListing)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Lowongan magang berhasil diperbarui",
		"listing": updatedListing,
	})
}
// Menghapus lowongan magang berdasarkan ID
func DeleteInternshipListingByID(c echo.Context) error {
	id := c.Param("id")

	if err := config.DB.Where("id = ?", id).Delete(&entity.Internship_Listing{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Gagal menghapus pendaftaran magang",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Pendaftaran magang berhasil dihapus",
	})
}
// Fungsi ini digunakan untuk memilih kandidat berdasarkan ID dan rentang nilai IPK (GPA)
func SelectCandidatesByGPAID(c echo.Context) error {
    // Mendapatkan ID 
    candidateID := c.Param("id")
	// Tetapkan nilai minimum dan maksimum IPK yang diinginkan
    minGPA := 3.5
    maxGPA := 4.0

    // Menggunakan GORM untuk mengambil kandidat yang memenuhi syarat dalam rentang IPK
    var candidates []entity.Internship_ApplicationForm
    if err := config.DB.Where("ID = ? AND GPA >= ? AND GPA <= ?", candidateID, minGPA, maxGPA).First(&candidates).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
    }

    // Loop melalui kandidat dan periksa nilai GPA serta status pembatalan
    for _, candidate := range candidates {
        if candidate.IsCanceled {
            // Jika formulir dibatalkan, ubah status menjadi "canceled"
            candidate.Status = constants.StatusCanceled
        } else if candidate.GPA >= minGPA && candidate.GPA <= maxGPA {
            // Jika memenuhi syarat dalam rentang IPK, ubah status menjadi "accepted"
            candidate.Status = constants.StatusAccepted
        } else {
            // Jika tidak memenuhi syarat dan tidak dibatalkan, ubah status menjadi "rejected"
            candidate.Status = constants.StatusRejected
        }

        if err := config.DB.Save(&candidate).Error; err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
        }

        // Buat entri di Selected_Candidate atau menyimpan data terkait kandidat yang dipilih 
        selectedCandidate := entity.Selected_Candidate{
            InternshipApplicationFormID: candidate.ID,
            InternshipApplicationForm:   candidate,
        }
        if err := config.DB.Create(&selectedCandidate).Error; err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
        }
    }

    return c.JSON(http.StatusOK, map[string]interface{}{"message": "Candidates selected based on GPA range"})
}
// Fungsi ini digunakan untuk mengirim email kepada kandidat yang diterima (status "accepted").
func SendEmailHandler(c echo.Context) error {
    userEmail := c.FormValue("userEmail")
    username := c.FormValue("username")
    status := c.FormValue("status")

    // Hanya kirim email jika status adalah "accepted" (dalam kasus Anda, "StatusAccepted")
    if status != constants.StatusAccepted {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email can only be sent for accepted candidates"})
    }

    err := helpers.SendEmailToUser(userEmail, username, status)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email", "details": err.Error()})
    }

    // Kirim respons sukses jika email terkirim
    return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully"})
}
// Fungsi ini digunakan untuk menampilkan semua kandidat yang ada di database.
func ViewAllCandidates(c echo.Context) error {
    // Mengambil semua kandidat dari database
    var candidates []entity.Internship_ApplicationForm
    if err := config.DB.Find(&candidates).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
    }


    // Menampilkan daftar kandidat
    return c.JSON(http.StatusOK, candidates)
}







