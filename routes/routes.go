package routes

import (
	"miniproject/controllers"
	"miniproject/internships/handler"
	"miniproject/internships/usecase"
	"miniproject/middleware"

	"github.com/labstack/echo/v4"
)

func InitmyRoutes() *echo.Echo {
	e := echo.New()
	internshipUsecase := usecase.NewInternshipApplicationUsecase()
	internshipHandler := handler.NewInternshipHandler(internshipUsecase)
	e.POST("/recommendation", internshipHandler.SubmitApplication)
	
	middleware.LogMiddleware(e)

	// Rute-rute admin
	adminGroup := e.Group("/admin")
	adminGroup.POST("/register", controllers.RegisterAdmin)
	adminGroup.POST("/login", controllers.LoginAdminController)
	adminGroup.GET("/:id", controllers.GetAdminByID, middleware.JWTMiddleware())
	adminGroup.PUT("/:id", controllers.UpdateAdminController, middleware.JWTMiddleware())
	// internship admin
	adminGroup.POST("/internship", controllers.CreateInternshipListing, middleware.JWTMiddleware())
	adminGroup.PUT("/internship/:id", controllers.UpdateInternshipListingByID, middleware.JWTMiddleware())
	adminGroup.DELETE("/internship/:id", controllers.DeleteInternshipListingByID, middleware.JWTMiddleware())
	adminGroup.GET("/selected-candidates/:id", controllers.SelectCandidatesByGPAID, middleware.JWTMiddleware())
	adminGroup.GET("/candidates", controllers.ViewAllCandidates, middleware.JWTMiddleware()) 
	adminGroup.POST("/email", controllers.SendEmailHandler, middleware.JWTMiddleware())

	// Route untuk User
	userGroup := e.Group("/users")
	userGroup.POST("/register", controllers.RegisterUser)
	userGroup.POST("/login", controllers.LoginUserController)
	userGroup.GET("/all", controllers.GetAllUsers, middleware.JWTMiddleware())
	userGroup.GET("/:id", controllers.GetUserByID, middleware.JWTMiddleware())
	userGroup.PUT("/:id", controllers.UpdateUserByID, middleware.JWTMiddleware())
	userGroup.DELETE("/:id", controllers.DeleteUser, middleware.JWTMiddleware())
	// internship user
	userGroup.GET("/internship-listings", controllers.GetInternshipListings, middleware.JWTMiddleware())
	userGroup.POST("/apply-for-internship", controllers.ApplyForInternship, middleware.JWTMiddleware())
	userGroup.DELETE("/apply-for-internship/:id", controllers.CancelApplication, middleware.JWTMiddleware())
	userGroup.GET("/candidates", controllers.ViewAllCandidates, middleware.JWTMiddleware())
	userGroup.GET("/Application-Status/:id", controllers.GetApplicationStatus, middleware.JWTMiddleware())
	return e
}
