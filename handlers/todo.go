package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"todo-list/config"
	"todo-list/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// extractUserID safely gets the UserID from the Gin context
func extractUserID(c *gin.Context) uint {
	// The AuthMiddleware ensures this key exists and is of type uint
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	// Should not happen if middleware is correctly applied
	return 0 
}

// PublicTodoResponse maps a models.Todo to a PublicTodo DTO
func PublicTodoResponse(todo models.Todo) models.PublicTodo {
	return models.PublicTodo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Done:        todo.Done,
	}
}

// CreateTodo handles POST /todos
func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Set the UserID from the authenticated user
	todo.UserID = extractUserID(c)

	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo item"})
		return
	}

	c.JSON(http.StatusCreated, PublicTodoResponse(todo))
}

// GetTodos handles GET /todos?page=1&limit=10&status=done
func GetTodos(c *gin.Context) {
	userID := extractUserID(c)

	// --- Pagination & Filtering ---
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	status := c.Query("status") // 'done', 'not_done'
	sort := c.DefaultQuery("sort", "created_at desc") // e.g., 'title asc', 'created_at desc'

	var todos []models.Todo
	var total int64
	
	// Start query, restricted by UserID
	query := config.DB.Model(&models.Todo{}).Where("user_id = ?", userID)

	// Filtering
	if status == "done" {
		query = query.Where("done = ?", true)
	} else if status == "not_done" {
		query = query.Where("done = ?", false)
	}

	// Count total items
	query.Count(&total)

	// Fetch paginated, filtered, and sorted data
	if err := query.Order(sort).Limit(limit).Offset(offset).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}

	// Prepare response data
	var responseData []models.PublicTodo
	for _, todo := range todos {
		responseData = append(responseData, PublicTodoResponse(todo))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responseData,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// UpdateTodo handles PUT /todos/:id
func UpdateTodo(c *gin.Context) {
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo ID"})
		return
	}

	userID := extractUserID(c)

	// 1. Check ownership and existence
	var todo models.Todo
	// Find the todo only if it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Responding with Forbidden (403) as requested for unauthorized access/update
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"}) 
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 2. Bind the incoming JSON data
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// 3. Update the fields
	// Note: We use .Model(&todo).Updates() to update only the fields present in the map
	if err := config.DB.Model(&todo).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo item"})
		return
	}

	// GORM Updates() might not reload the model, so fetch the latest or return the updated struct
	config.DB.First(&todo, todo.ID)

	c.JSON(http.StatusOK, PublicTodoResponse(todo))
}

// DeleteTodo handles DELETE /todos/:id
func DeleteTodo(c *gin.Context) {
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo ID"})
		return
	}
	
	userID := extractUserID(c)
	
	// 1. Check ownership and existence before deleting
	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Responding with Forbidden (403) as requested for unauthorized access/delete
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"}) 
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 2. Delete the item
	if err := config.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete todo item"})
		return
	}

	// Respond with Status 204 No Content for successful deletion
	c.Status(http.StatusNoContent)
}