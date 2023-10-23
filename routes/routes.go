package routes

import (
	"miniproject/delivery/controllers"

	"github.com/labstack/echo"
)

func InitmyRoutes(A *controllers.AdminController, u *controllers.UserController) *echo.Echo {
	e := echo.New()

	// Rute-rute admin
	adminGroup := e.Group("/admin")
	adminGroup.POST("/login", A.LoginAdmin)
	adminGroup.GET("/admin/:id", A.GetAdminByID)
	adminGroup.PUT("/admin/:id", A.UpdateAdmin)
	adminGroup.DELETE("/admin/:id", A.DeleteAdmin)
	adminGroup.POST("/internship", A.CreateInternshipListing)
	adminGroup.POST("/application", A.CreateInternshipApplicationForm)
	adminGroup.PUT("/updateStatus/:id", A.UpdateApplicationStatus)
	adminGroup.PUT("/admin/verifyCancel/:id", A.VerifyCancelApplication)

	// Route untuk registrasi
	userGroup := e.Group("/user")
	userGroup.POST("/register", u.Register)
	userGroup.POST("/login", u.Login)
	userGroup.PUT("/profile", u.UpdateProfile)
	userGroup.POST("/profilePicture", u.UploadProfilePicture)
	userGroup.GET("/users", u.GetAllUsers)
	userGroup.GET("/users/:id", u.GetUserByID)
	userGroup.DELETE("/users/:id", u.DeleteUserByID)
	userGroup.GET("/internshipListings", u.GetInternshipListings)
	userGroup.POST("/internship/:id", u.ChooseInternshipListing)
	userGroup.GET("/ApplicationStatus", u.GetApplicationStatus)
	userGroup.POST("/cancelApplication/:id", u.CancelApplication)

	return e
}
