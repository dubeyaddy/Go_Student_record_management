package routes

import (
	"student-record-management/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/students", controllers.CreateStudent)
	r.GET("/students", controllers.GetStudents)
	r.GET("/students/:id", controllers.GetStudentByID)
	r.GET("/students/type/:studentTypeId", controllers.GetStudentByStudentType)
	r.PUT("/students/:id", controllers.UpdateStudent)
	r.DELETE("/students/:id", controllers.DeleteStudent)

	r.POST("/studentTypes", controllers.CreateStudentType)
	r.GET("/studentTypes", controllers.GetStudentTypes)
	r.GET("/studentTypes/:id", controllers.GetStudentTypeByID)
	r.PUT("/studentTypes/:id", controllers.UpdateStudentType)
	r.DELETE("/studentTypes/:id", controllers.DeleteStudentType)

	return r
}
