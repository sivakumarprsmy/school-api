package models

import (
	"time"
)

type Student struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null"                 json:"name"`
	Email     string    `gorm:"uniqueIndex;not null"     json:"email"`
	Age       int       `gorm:"not null"                 json:"age"`
	Grade     string    `gorm:"not null"                 json:"grade"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStudentRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"   binding:"required,min=1,max=120"`
	Grade string `json:"grade" binding:"required"`
}
