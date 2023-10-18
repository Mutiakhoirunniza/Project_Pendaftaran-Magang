package main

import (
	App "project/App/config"
	"project/middleware"
	"project/routes"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	database.InitDBMysql()

	e.Use(middleware.InitLoggerMiddleware())

	// Inisialisasi rute-rute
	routes.InitmyRoutes(e, db)

	// Jalankan server Echo
	e.Start(":8080")
}
