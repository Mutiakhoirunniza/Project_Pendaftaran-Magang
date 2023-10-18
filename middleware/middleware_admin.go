package middleware

import (
    "github.com/labstack/echo"
    "net/http"
)

func AdminMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Dapatkan peran pengguna/admin dari konteks
            userRole, ok := c.Get("userRole").(string)
            if !ok {
                return echo.NewHTTPError(http.StatusUnauthorized, "Token JWT tidak valid")
            }

            // Periksa apakah pengguna adalah admin berdasarkan perannya
            if userRole != "admin" {
                return echo.NewHTTPError(http.StatusForbidden, "Anda tidak memiliki izin admin")
            }

            return next(c)
        }
    }
}

// Middleware untuk memastikan hak admin
func PastikanPenggunaAdmin(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        isAdmin, ok := c.Get("isAdmin").(bool)
        if !ok || !isAdmin {
            return echo.NewHTTPError(http.StatusUnauthorized, "Anda tidak memiliki izin admin")
        }
        return next(c)
    }
}
