package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"student-record-management/config"
	"student-record-management/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStudentType(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	body := []byte(`{"type":"School"}`)
	req, _ := http.NewRequest("POST", "/studentTypes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "School")
}

func TestCreateStudentType_Duplicate(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Seed one type
	config.DB.Create(&model.StudentType{Type: "College"})

	body := []byte(`{"type":"College"}`)
	req, _ := http.NewRequest("POST", "/studentTypes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
}

func TestGetStudentTypes(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	config.DB.Create(&model.StudentType{Type: "University"})

	req, _ := http.NewRequest("GET", "/studentTypes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "University")
}

func TestGetStudentTypeByID(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	st := model.StudentType{Type: "School"}
	config.DB.Create(&st)

	req, _ := http.NewRequest("GET", "/studentTypes/"+fmt.Sprint(st.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "School")
}

func TestUpdateStudentType(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	st := model.StudentType{Type: "OldType"}
	config.DB.Create(&st)

	body := []byte(`{"type":"NewType"}`)
	req, _ := http.NewRequest("PUT", "/studentTypes/"+fmt.Sprint(st.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "NewType")
}

func TestDeleteStudentType(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	st := model.StudentType{Type: "ToDelete"}
	config.DB.Create(&st)

	req, _ := http.NewRequest("DELETE", "/studentTypes/"+fmt.Sprint(st.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "deleted successfully")

	// Verify deletion
	var check model.StudentType
	err := config.DB.First(&check, st.ID).Error
	assert.NotNil(t, err)
}
