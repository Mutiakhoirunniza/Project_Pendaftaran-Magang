package routes

import (
	"project/controllers"
	"project/middleware"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func InitmyRoutes(e *echo.Echo, db *gorm.DB) {
	userController := controllers.NewUserController(db)
	adminController := controllers.NewInternshipController(db)

	// Grup untuk rute pengguna (user)
	userGroup := e.Group("/user")
	userGroup.POST("/register", userController.RegisterUser)
	userGroup.POST("/login", userController.LoginUser)
	userGroup.GET("/profile", userController.GetUserProfile)
	userGroup.GET("/profile/:id", userController.GetUserProfileByID)
	userGroup.PUT("/profile", userController.UpdateUserProfile)
	userGroup.DELETE("/profile/:id", userController.DeleteUser)
	userGroup.GET("/application-status", userController.GetApplicationStatus)
	userGroup.DELETE("/application/:id", userController.CancelApplication)

	// Grup untuk rute admin
    adminGroup := e.Group("/admin")
    adminGroup.Use(middleware.PastikanPenggunaAdmin)
    adminGroup.POST("/AdminInternship", adminController.CreateInternship)
    adminGroup.PUT("/AdminInternship/:id", adminController.UpdateInternship)
    adminGroup.DELETE("/AdminInternship/:id", adminController.DeleteInternship)
    adminGroup.PUT("/AdminApplicationStatus/:id", adminController.UpdateApplicationStatus)
    adminGroup.POST("/login", adminController.AdminLogin) 
	adminGroup.GET("/AdminApplicationStatus", adminController.ApplicationsByStatus) 
	
}
