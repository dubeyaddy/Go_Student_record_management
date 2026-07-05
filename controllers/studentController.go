package controllers

import (
	"net/http"
	"student-record-management/config"
	"student-record-management/model"

	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	var student model.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var studentType model.StudentType

	result := config.DB.First(&studentType, student.StudentTypeId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "studentTypeId not found",
		})
		return
	}

	config.DB.Create(&student)

	c.JSON(http.StatusOK, student)
}

func GetStudents(c *gin.Context) {
	var students []model.Student
	config.DB.Find(&students)
	c.JSON(http.StatusOK, students)
}

func GetStudentByID(c *gin.Context) {
	id := c.Param("id")
	var student model.Student

	if err := config.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func GetStudentByStudentType(c *gin.Context) {
	studentTypeId := c.Param("studentTypeId")
	var student model.Student

	if err := config.DB.Where("student_type_id = ?", studentTypeId).Find(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var student model.Student

	// Check if student exists
	if err := config.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Bind new data
	var updatedData model.Student
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	student.Name = updatedData.Name
	student.StudentTypeId = updatedData.StudentTypeId
	student.GuardianName = updatedData.GuardianName

	config.DB.Save(&student)
	c.JSON(http.StatusOK, student)
}

func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	var student model.Student

	// Check if student exists
	if err := config.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Delete student
	config.DB.Delete(&student)
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
