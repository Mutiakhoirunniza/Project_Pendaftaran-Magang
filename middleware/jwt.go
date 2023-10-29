package middleware

import (
	"errors"
	"miniproject/constants"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

//	func JWTMiddleware() echo.MiddlewareFunc {
//		return echojwt.WithConfig(echojwt.Config{
//			SigningKey:    []byte(constants.SecretKey),
//			SigningMethod: "HS256",
//		})
//	}
func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(constants.SecretKey),
	})
}

func CreateToken(username string, role string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 24 * 7).Unix(),
			Issuer:    "miniproject",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(constants.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenAdmin(c echo.Context) (string, error) {
	adminToken := c.Get("admin")
	if adminToken == nil {
		return "", errors.New("Admin token not found in context")
	}
	token, ok := adminToken.(*jwt.Token)
	if !ok {
		return "", errors.New("Invalid admin token in context")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", errors.New("Invalid claims in token")
	}

	if claims.Role != "admin" {
		return "", errors.New("Not an admin")
	}

	return claims.Username, nil
}

func ExtractTokenUser(c echo.Context) (string, error) {
	UserToken := c.Get("user")
	if UserToken == nil {
		return "", errors.New("User token not found in context")
	}
	token, ok := UserToken.(*jwt.Token)
	if !ok {
		return "", errors.New("Invalid User token in context")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", errors.New("Invalid claims in token")
	}

	if claims.Role != "user" {
		return "", errors.New("Not an user")
	}
	return claims.Username, nil
}
