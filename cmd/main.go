package main

import (
	"log"
	"miniproject/delivery/controllers"
	"miniproject/routes"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load the .env file")
	}

	// Initialize the database connection
	dsn := os.Getenv("DBUSER") + ":" + os.Getenv("DBPASS") + "@tcp(" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + ")/" + os.Getenv("DBNAME") + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	// Initialize logger
	logger := log.New(os.Stdout, "MYAPP: ", log.Ldate|log.Ltime)

	// Create instances of your controllers
	adminController := controllers.NewAdminController(db, logger)
	userController := controllers.NewUserController(db, logger)

	server := routes.InitmyRoutes(adminController, userController)

	// Start the server on a separate line
	server.Start(":" + os.Getenv("SERVERPORT"))
}
