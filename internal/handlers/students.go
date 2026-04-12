package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/periasamy/school-api/internal/models"
)

type StudentHandler struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewStudentHandler(db *gorm.DB, log zerolog.Logger) *StudentHandler {
	return &StudentHandler{db: db, log: log}
}

// CreateStudent handles POST /api/v1/students
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req models.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := models.Student{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Grade: req.Grade,
	}

	if err := h.db.WithContext(c.Request.Context()).Create(&student).Error; err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "student with this email already exists"})
			return
		}
		h.log.Error().Err(err).Msg("failed to create student")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// GetStudentByID handles GET /api/v1/students/:id
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var student models.Student
	if err := h.db.WithContext(c.Request.Context()).First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		h.log.Error().Err(err).Uint64("id", id).Msg("failed to fetch student")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// GetAllStudents handles GET /api/v1/students
func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	var students []models.Student

	if err := h.db.WithContext(c.Request.Context()).Find(&students).Error; err != nil {
		h.log.Error().Err(err).Msg("failed to fetch students")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

// DeleteStudent handles DELETE /api/v1/students/:id
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	result := h.db.WithContext(c.Request.Context()).Delete(&models.Student{}, id)
	if result.Error != nil {
		h.log.Error().Err(result.Error).Uint64("id", id).Msg("failed to delete student")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint")
}
