package middleware

import (
	"errors"
	"fmt"
	"miniproject/constants"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func GenerateJWTToken(username string) (string, error) {
	// Menyiapkan claim (klaim) token
	claims := jwt.MapClaims{
		"username": username,
		"exp":      0, // Token tidak memiliki masa berlaku
	}

	// Membuat token dengan claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key
	secretKey := []byte(constants.SecretKey)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JWTMiddleware adalah middleware Echo untuk otentikasi dengan JWT
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak ditemukan")
			}

			// Parse token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Pastikan Anda menggunakan kunci yang sesuai dengan yang digunakan saat membuat token
				return []byte(constants.SecretKey), nil
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
			}

			// Mengekstrak klaim (claims) dari token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak mengandung klaim yang valid")
			}

			// Mengekstrak ID pengguna dari klaim (asumsi ID disimpan dalam klaim "sub")
			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Tidak dapat menemukan ID pengguna dalam token JWT")
			}

			// Set userID dalam konteks Echo sehingga dapat diakses di handler selanjutnya
			c.Set("userID", userID)
			return next(c)
		}
	}
}

func GetUserIDFromToken(c echo.Context) (string, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return "", errors.New("Token tidak ada")
	}

	claims, err := verifyJWTToken(token)
	if err != nil {
		return "", err
	}

	// Ambil ID pengguna dari klaim token JWT
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("Token tidak berisi ID pengguna")
	}

	return userID, nil
}

func verifyJWTToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte(constants.SecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metode penandatanganan token tidak valid")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Periksa apakah token telah kedaluwarsa
		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("Token JWT telah kedaluwarsa")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("Token JWT tidak valid")
}

// // Middleware untuk memeriksa apakah pengguna telah mengisi formulir pendaftaran
// func CheckInternshipFormSubmitted() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			// Cek apakah pengguna telah mengisi formulir
// 			var formData entity.InternshipApplicationForm
// 			if err := u.DB.Where("user_id = ?", c.Get("user_id")).First(&formData).Error; err != nil {
// 				// Pengguna belum mengisi formulir, kembalikan pesan error
// 				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Anda harus mengisi formulir pendaftaran terlebih dahulu"})
// 			}

// 			// Lanjutkan ke tahap memilih daftar lowongan magang
// 			return next(c)
// 		}
// 	}
// }

// // user
// func DatabaseMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		c.Set("db", database.DB) // Mengatur koneksi database ke konteks Echo
// 		return next(c)
// 	}
// }

// func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		user := c.Get("user").(*entity.User)
// 		if user == nil || !user.IsAdmin {
// 			return c.JSON(http.StatusUnauthorized, "Akses ditolak: Hanya admin yang dapat mengakses halaman ini.")
// 		}
// 		return next(c)
// 	}
// }
