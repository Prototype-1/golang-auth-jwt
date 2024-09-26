package routes

import (
	"github.com/Prototype-1/golang-auth-jwt/controllers"
	"github.com/Prototype-1/golang-auth-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.POST("/signup", controllers.AdminSignup)
		admin.POST("/login", controllers.AdminLogIn)

		// Protect routes with authentication middleware
		admin.Use(middleware.AuthMiddleware())

		admin.GET("/get_users", middleware.AuthRequired(), controllers.GetAllUsers)
		admin.POST("/create_user", controllers.CreateUser)
		admin.GET("/users/:user_id", controllers.GetUser)
		admin.PUT("/update_user/:id", controllers.UpdateUser)
		admin.DELETE("/delete_user/:id", controllers.DeleteUser)
		
		admin.StaticFile("/dashboard", "./frontend/dashboard.html")
	}
}




