package routes

import (
	"miniproject/controllers"

	"github.com/labstack/echo/v4"
)

func InitmyRoutes(e *echo.Echo) {
	// e := echo.New()

	// Rute-rute admin
	adminGroup := e.Group("/admin")
	adminGroup.POST("/login", controllers.LoginAdminController)
	adminGroup.GET("/admin/:id", controllers.GetAdminByID)
	adminGroup.PUT("/admin/:id", controllers.UpdateAdmin)
	adminGroup.DELETE("/admin/:id", controllers.DeleteAdmin)
	adminGroup.POST("/internship", controllers.CreateInternshipListing)
	adminGroup.POST("/application", controllers.CreateInternshipApplicationForm)
	adminGroup.PUT("/updateStatus/:id", controllers.UpdateApplicationStatus)
	adminGroup.PUT("/admin/verifyCancel/:id", controllers.VerifyCancelApplication)

	// Route untuk registrasi
	userGroup := e.Group("/user")
	userGroup.POST("/register", controllers.Register)
	userGroup.POST("/login", controllers.LoginUserController)
	userGroup.PUT("/profilePicture", controllers.UpdateProfileAndUploadPicture)
	userGroup.GET("/users", controllers.GetAllUsers)
	userGroup.GET("/users/:id", controllers.GetUserByID)
	userGroup.DELETE("/users/:id", controllers.DeleteUserByID)
	userGroup.GET("/internshipListings", controllers.GetInternshipListings)
	userGroup.POST("/internship/:id", controllers.ChooseInternshipListing)
	userGroup.GET("/ApplicationStatus", controllers.GetApplicationStatus)
	userGroup.POST("/cancelApplication/:id", controllers.CancelApplication)
}
