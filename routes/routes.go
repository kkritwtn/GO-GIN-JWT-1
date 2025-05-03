package routes

import (
	adminControllers "go-jwt-crud/controllers/admin"
	userControllers "go-jwt-crud/controllers/user"
	"go-jwt-crud/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Admin routes
	admin := r.Group("/api/admin")
	{
		admin.POST("/register", adminControllers.RegisterAdmin)
		admin.POST("/login", adminControllers.LoginAdmin)
	}

	// Protected Admin routes
	adminAuth := r.Group("/api/admin")
	adminAuth.Use(middleware.AdminAuthMiddleware()) // Middleware that checks role == "admin"
	{
		adminAuth.GET("/dashboard", adminControllers.AdminDashboard)
		adminAuth.GET("/users", adminControllers.GetAllUsers)
		// adminAuth.DELETE("/user/:id", adminControllers.DeleteUser)
	}

	// User public routes
	user := r.Group("/api/user")
	{
		user.POST("/register", userControllers.RegisterUser)
		user.POST("/login", userControllers.LoginUser)
	}

	// Protected User routes
	userAuth := r.Group("/api/user")
	userAuth.Use(middleware.UserAuthMiddleware()) // Middleware that checks role == "user"
	{
		userAuth.GET("/profile", userControllers.GetProfile)
		userAuth.PUT("/update", userControllers.UpdateProfile)
		userAuth.DELETE("/delete", userControllers.DeleteAccount)
	}
}
