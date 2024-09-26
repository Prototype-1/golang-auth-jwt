package routes

import (
	"github.com/Prototype-1/golang-auth-jwt/controllers"
	"github.com/Prototype-1/golang-auth-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/users")
	{
		user.POST("/signup", controllers.SignUp)
		user.POST("/login", controllers.Login)

		// Protect routes with authentication middleware
		user.Use(middleware.AuthMiddleware())
		user.GET("/me", controllers.GetAllUsers)
		user.StaticFile("/home", "./frontend/user_home.html")
	}
}








