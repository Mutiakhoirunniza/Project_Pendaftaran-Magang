package middleware

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	godotenv.Load(".env")
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(os.Getenv("SecretKey")),
		SigningMethod: "HS256",
	})
}

func CreateToken(userId uint, username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SecretKey")))
}

func ExtractToken(e echo.Context) (uint, string) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userId"].(float64)
		username := claims["username"].(string)
		return uint(userId), username
	}
	return 0, ""
}
