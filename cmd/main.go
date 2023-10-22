package main

import (
	"miniproject/delivery/controllers"
	"miniproject/routes"
)

func main() {

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// Inisialisasi controllers
	adminControllers := controllers.NewAdminController(db)
	internshipControllers := controllers.NewInternshipController(db)

	// Set up router
	e := routes.Initmyrouter(adminControllers, internshipControllers)
	e.Logger.Fatal(e.Start(":8080"))

}
