package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"student-record-management/config"
	"student-record-management/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup test DB
func setupTestDB() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	config.DB.AutoMigrate(&model.Student{}, &model.StudentType{})
}

// Setup router with routes
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/students", CreateStudent)
	r.GET("/students", GetStudents)
	r.GET("/students/:id", GetStudentByID)
	r.PUT("/students/:id", UpdateStudent)
	r.DELETE("/students/:id", DeleteStudent)
	r.POST("/studentTypes", CreateStudentType)
	r.GET("/studentTypes", GetStudentTypes)
	r.GET("/studentTypes/:id", GetStudentTypeByID)
	r.PUT("/studentTypes/:id", UpdateStudentType)
	r.DELETE("/studentTypes/:id", DeleteStudentType)
	return r
}

func TestCreateStudent(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// First create a StudentType so foreign key exists
	studentType := model.StudentType{Type: "Primary"}
	config.DB.Create(&studentType)

	// Prepare request body
	student := model.Student{
		Name:          "John Doe",
		StudentTypeId: studentType.ID,
		GuardianName:  "Jane Doe",
	}
	body, _ := json.Marshal(student)

	req, _ := http.NewRequest("POST", "/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var created model.Student
	json.Unmarshal(w.Body.Bytes(), &created)
	assert.Equal(t, "John Doe", created.Name)
	assert.Equal(t, studentType.ID, created.StudentTypeId)
}

func TestGetStudents(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Seed data
	studentType := model.StudentType{Type: "Secondary"}
	config.DB.Create(&studentType)
	student := model.Student{Name: "Alice", StudentTypeId: studentType.ID, GuardianName: "Bob"}
	config.DB.Create(&student)

	req, _ := http.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Alice")
}

func TestGetStudentByID_NotFound(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/students/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Student not found")
}

func TestUpdateStudent(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Seed data
	studentType := model.StudentType{Type: "Primary"}
	config.DB.Create(&studentType)
	student := model.Student{Name: "Old Name", StudentTypeId: studentType.ID, GuardianName: "Old Guardian"}
	config.DB.Create(&student)

	// Prepare updated data
	updated := model.Student{Name: "New Name", StudentTypeId: studentType.ID, GuardianName: "New Guardian"}
	body, _ := json.Marshal(updated)

	req, _ := http.NewRequest("PUT", "/students/"+fmt.Sprint(student.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "New Name")
}

func TestDeleteStudent(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Seed data
	studentType := model.StudentType{Type: "Secondary"}
	config.DB.Create(&studentType)
	student := model.Student{Name: "ToDelete", StudentTypeId: studentType.ID, GuardianName: "Guardian"}
	config.DB.Create(&student)

	req, _ := http.NewRequest("DELETE", "/students/"+fmt.Sprint(student.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "deleted successfully")

	// Verify deletion
	var check model.Student
	err := config.DB.First(&check, student.ID).Error
	assert.NotNil(t, err) // should not find student
}

func TestGetStudentByStudentType(t *testing.T) {
	setupTestDB()
	router := gin.Default()
	router.GET("/students/type/:studentTypeId", GetStudentByStudentType)

	// Seed data
	studentType := model.StudentType{Type: "Primary"}
	config.DB.Create(&studentType)
	student := model.Student{Name: "ByType", StudentTypeId: studentType.ID, GuardianName: "Guardian"}
	config.DB.Create(&student)

	req, _ := http.NewRequest("GET", "/students/type/"+fmt.Sprint(studentType.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ByType")
}
