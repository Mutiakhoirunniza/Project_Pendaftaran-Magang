package middleware

import (
	"net/http"
	constants "project/contants"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak ditemukan")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(constants.SecretKey), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak mengandung klaim yang valid")
			}

			// Mengekstrak ID pengguna dari klaim (asumsi ID disimpan dalam klaim "sub")
			userID, ok := claims["sub"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Tidak dapat menemukan ID pengguna dalam token JWT")
			}

			// Simpan ID pengguna ke konteks agar dapat diakses oleh penanganan permintaan selanjutnya
			c.Set("userID", userID)

			return next(c)
		}
	}
}
