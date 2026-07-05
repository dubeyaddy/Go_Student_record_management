package main

import (
	"log"
	"os"
	"student-record-management/config"
	"student-record-management/controllers"
	"student-record-management/model"
	"student-record-management/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDatabase()
	config.DB.AutoMigrate(&model.Student{}, &model.StudentType{})
	controllers.SeedStudentTypes()

	r := routes.SetupRouter()
	r.Run(":" + os.Getenv("PORT"))
}
