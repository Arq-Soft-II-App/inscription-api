package controllers

import (
	"inscription-api/src/dto"
	"inscription-api/src/errors"
	"inscription-api/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InscriptionController struct {
	service services.InscriptionService
	logger  *zap.Logger
}

func NewInscriptionController(service services.InscriptionService, logger *zap.Logger) *InscriptionController {
	return &InscriptionController{service: service, logger: logger}
}

func (ic *InscriptionController) EnrollStudent(c *gin.Context) {
	var enrollRequest dto.EnrollRequestResponseDto
	if err := c.ShouldBindJSON(&enrollRequest); err != nil {
		appErr := errors.ErrInvalidData
		ic.logger.Warn("[INSCRIPTION-API] Datos inválidos en EnrollStudent", zap.Error(err))
		c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
		return
	}

	if err := ic.service.EnrollStudent(c.Request.Context(), enrollRequest.CourseId, enrollRequest.UserId); err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServer.Error()})
		return
	}

	ic.logger.Info("[INSCRIPTION-API] Estudiante inscrito exitosamente", zap.String("user_id", enrollRequest.UserId), zap.String("course_id", enrollRequest.CourseId))
	c.Status(http.StatusCreated)
}

func (ic *InscriptionController) GetMyCourses(c *gin.Context) {
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		appErr := errors.ErrMissingUserId
		ic.logger.Warn("[INSCRIPTION-API] UserId es requerido en GetMyCourses")
		c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
		return
	}

	courses, err := ic.service.GetMyCourses(c.Request.Context(), userIdStr)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			ic.logger.Error("[INSCRIPTION-API] Error al obtener los cursos", zap.Error(appErr))
			c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
			return
		}
		ic.logger.Error("[INSCRIPTION-API] Error interno al obtener los cursos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServer.Error()})
		return
	}

	ic.logger.Info("[INSCRIPTION-API] Cursos devueltos para el usuario", zap.String("user_id", userIdStr), zap.Int("total_courses", len(courses)))
	c.JSON(http.StatusOK, courses)
}

func (ic *InscriptionController) GetStudentsInCourse(c *gin.Context) {
	courseIdStr := c.Param("cid")
	if courseIdStr == "" {
		appErr := errors.ErrMissingCourseId
		ic.logger.Warn("[INSCRIPTION-API] CourseId es requerido en GetStudentsInCourse")
		c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
		return
	}

	students, err := ic.service.GetStudentsInCourse(c.Request.Context(), courseIdStr)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			ic.logger.Error("[INSCRIPTION-API] Error al obtener los estudiantes", zap.Error(appErr))
			c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
			return
		}
		ic.logger.Error("[INSCRIPTION-API] Error interno al obtener los estudiantes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServer.Error()})
		return
	}

	ic.logger.Info("[INSCRIPTION-API] Estudiantes devueltos para el curso", zap.String("course_id", courseIdStr), zap.Int("total_students", len(students)))
	c.JSON(http.StatusOK, students)
}

func (ic *InscriptionController) IsEnrolled(c *gin.Context) {
	courseIdStr := c.Param("cid")
	userIdStr := c.Param("userId")

	if courseIdStr == "" || userIdStr == "" {
		ic.logger.Warn("[INSCRIPTION-API] CourseId y UserId son requeridos en IsEnrolled")
		c.JSON(http.StatusBadRequest, gin.H{"error": "CourseId y UserId son requeridos"})
		return
	}

	isEnrolled, err := ic.service.IsEnrolled(c.Request.Context(), courseIdStr, userIdStr)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			ic.logger.Error("[INSCRIPTION-API] Error al verificar la inscripción", zap.Error(appErr))
			c.JSON(appErr.HTTPStatusCode, gin.H{"error": appErr.Error()})
			return
		}
		ic.logger.Error("[INSCRIPTION-API] Error interno al verificar la inscripción", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServer.Error()})
		return
	}

	ic.logger.Info("[INSCRIPTION-API] Resultado de IsEnrolled", zap.String("user_id", userIdStr), zap.String("course_id", courseIdStr), zap.Bool("is_enrolled", isEnrolled))
	c.JSON(http.StatusOK, gin.H{"enrolled": isEnrolled})
}
