// routes/inscription_routes.go
package routes

import (
	"inscription-api/src/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, inscriptionController *controllers.InscriptionController) {
	inscriptionRoutes := router.Group("/")
	{
		inscriptionRoutes.POST("/enroll", inscriptionController.EnrollStudent)
		inscriptionRoutes.GET("/myCourses", inscriptionController.GetMyCourses)
		inscriptionRoutes.GET("/studentsInThisCourse/:cid", inscriptionController.GetStudentsInCourse)
		inscriptionRoutes.GET("/isEnrolled/:cid/:userId", inscriptionController.IsEnrolled)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
}
