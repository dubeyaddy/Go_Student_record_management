package controllers

import (
	"net/http"
	"student-record-management/config"
	"student-record-management/model"

	"github.com/gin-gonic/gin"
)

func SeedStudentTypes() {
	var count int64

	config.DB.Model(&model.StudentType{}).Count(&count)

	if count == 0 {
		studentTypes := []model.StudentType{
			{ID: 1, Type: "School"},
			{ID: 2, Type: "College"},
			{ID: 3, Type: "University"},
		}

		config.DB.Create(&studentTypes)
	}
}

// Create
func CreateStudentType(c *gin.Context) {
	var studentType model.StudentType
	if err := c.ShouldBindJSON(&studentType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if type already exists
	var existing model.StudentType
	if err := config.DB.Where("type = ?", studentType.Type).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "StudentType already exists"})
		return
	}

	config.DB.Create(&studentType)
	c.JSON(http.StatusOK, studentType)
}

// Get all
func GetStudentTypes(c *gin.Context) {
	var studentTypes []model.StudentType
	config.DB.Find(&studentTypes)
	c.JSON(http.StatusOK, studentTypes)
}

// Get by ID
func GetStudentTypeByID(c *gin.Context) {
	id := c.Param("id")
	var studentType model.StudentType
	if err := config.DB.First(&studentType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "StudentType not found"})
		return
	}
	c.JSON(http.StatusOK, studentType)
}

// Update
func UpdateStudentType(c *gin.Context) {
	id := c.Param("id")
	var studentType model.StudentType

	if err := config.DB.First(&studentType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "StudentType not found"})
		return
	}

	var updatedData model.StudentType
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentType.Type = updatedData.Type
	config.DB.Save(&studentType)
	c.JSON(http.StatusOK, studentType)
}

// Delete
func DeleteStudentType(c *gin.Context) {
	id := c.Param("id")
	var studentType model.StudentType

	if err := config.DB.First(&studentType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "StudentType not found"})
		return
	}

	config.DB.Delete(&studentType)
	c.JSON(http.StatusOK, gin.H{"message": "StudentType deleted successfully"})
}
