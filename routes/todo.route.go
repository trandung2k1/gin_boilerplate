package routes

import (
	"ginboilerplate/controllers"
	"ginboilerplate/middleware"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(rg *gin.RouterGroup) {
	todos := rg.Group("/todos")
	todos.Use(middleware.Logger())
	todos.GET("/", controllers.GetAllTodo)
	todos.POST("/", middleware.VerifyToken(), controllers.CreateTodo)
	todos.GET("/:id", middleware.VerifyToken(), controllers.GetTodoById)
}
