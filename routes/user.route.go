package routes

import (
	"ginboilerplate/controllers"
	"ginboilerplate/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/auth")
	users.POST("/register", controllers.RegisterAccount)
	users.POST("/login", controllers.Login)
	users.GET("/users/:id", middleware.VerifyToken(), controllers.GetUserById)
}
