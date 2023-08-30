package main

import (
	_ "assignment_1/docs" // This is required for Swagger to find your docs
	"assignment_1/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           TODO Application
// @version         1.0
// @description     This is a todo list management application.
// @termsOfService  http://swagger.io/terms/
func main() {
	r := gin.Default()

	// Swagger UIgog
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/todos", handlers.GetTodos)
		v1.GET("/todos/:id", handlers.GetTodoByID)
		v1.POST("/todos", handlers.CreateTodo)
		v1.PUT("/todos/:id", handlers.UpdateTodo)
		v1.DELETE("/todos/:id", handlers.DeleteTodo)
	}

	r.Run(":8080")
}
