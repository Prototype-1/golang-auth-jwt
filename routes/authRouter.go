package routes

import (
	"github.com/Prototype-1/golang-auth-jwt/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUp)  // Ensure this is not duplicated
		auth.POST("/login", controllers.Login)    // Ensure this is not duplicated
	}
}




