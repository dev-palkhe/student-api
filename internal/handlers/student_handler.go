package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dev-palkhe/student-api/internal/models"
	"github.com/dev-palkhe/student-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentHandler struct {
	repo repository.StudentRepository
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student.ID = uuid.New()

	createdStudent, err := h.repo.Create(c.Request.Context(), student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, createdStudent)
}

func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	students, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get students"})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	student, err := h.repo.GetByID(c.Request.Context(), id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student.ID = id
	student.UpdatedAt = time.Now()

	err = h.repo.Update(c.Request.Context(), student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}
	updatedStudent, err := h.repo.GetByID(c.Request.Context(), student.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated student"})
		return
	}

	c.JSON(http.StatusOK, updatedStudent)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
}
