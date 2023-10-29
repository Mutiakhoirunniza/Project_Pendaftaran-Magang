package main

import (
	"miniproject/infra/config"
	"miniproject/infra/database"
	"miniproject/infra/migration"
	"miniproject/routes"
	"os"
)

func main() {
	// Inisialisasi konfigurasi
	cfg := config.InitConfig()
	db := database.InitDBMysql(cfg)
	migration.InitMigrationMysql(db)
	e := routes.InitmyRoutes()
	
	

	// Mulai server
	e.Logger.Fatal(e.Start(":" + os.Getenv("SERVERPORT")))
}

// Start(":" + os.Getenv("SERVERPORT"))
// logger := log.New(os.Stdout, "MYAPP: ", log.Ldate|log.Ltime)

// // Create instances of your controllers
// adminController := controllers.NewAdminController(db, logger)
// userController := controllers.NewUserController(db, logger)

// routes.InitmyRoutes(adminController, userController)

// // Load environment variables from the .env file
// err := godotenv.Load()
// if err != nil {
// 	log.Fatal("Failed to load the .env file")
// }

// // Initialize the database connection
// dsn := os.Getenv("DBUSER") + ":" + os.Getenv("DBPASS") + "@tcp(" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + ")/" + os.Getenv("DBNAME") + "?parseTime=true"
// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// if err != nil {
// 	log.Fatal("Failed to connect to the database")
// }
