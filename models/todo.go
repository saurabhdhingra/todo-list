package models

import "gorm.io/gorm"

// Todo represents a single to-do item
type Todo struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done" gorm:"default:false"`
	UserID      uint   `json:"user_id"` // Foreign key
	User        User   // GORM association field
}

// PublicTodo is a DTO used for list and single item responses
type PublicTodo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}