package main

import (
	"log"
	"net/http"

	"todo-list/config"
	"todo-list/handlers"
	"todo-list/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize Database and Auto-Migrate Schemas
	config.ConnectDatabase()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 2. Public Routes (No Authentication required)
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)

	// 3. Authenticated Routes Group
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware()) // Apply JWT Middleware to this group
	{
		// Todo CRUD Endpoints
		authenticated.POST("/todos", handlers.CreateTodo)
		authenticated.GET("/todos", handlers.GetTodos)
		authenticated.PUT("/todos/:id", handlers.UpdateTodo)
		authenticated.DELETE("/todos/:id", handlers.DeleteTodo)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "message": "Todo API is running"})
	})

	log.Println("Server listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}