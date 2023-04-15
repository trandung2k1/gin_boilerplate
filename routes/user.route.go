package routes

import (
	"ginboilerplate/controllers"
	"ginboilerplate/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/auth")
	users.GET("/users/:id", middleware.VerifyToken(), controllers.GetUserById)
}
