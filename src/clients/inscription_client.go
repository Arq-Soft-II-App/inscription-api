package clients

import (
	"context"
	"inscription-api/src/errors"
	"inscription-api/src/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InscriptionRepository interface {
	Create(ctx context.Context, inscription *models.Inscripto) error
	GetMyCourses(ctx context.Context, userId string) ([]string, error)
	GetStudentsInCourse(ctx context.Context, courseId string) ([]string, error)
	IsEnrolled(ctx context.Context, courseId string, userId string) (bool, error)
}

type gormInscriptionRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGormInscriptionRepository(db *gorm.DB, logger *zap.Logger) InscriptionRepository {
	return &gormInscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormInscriptionRepository) Create(ctx context.Context, inscription *models.Inscripto) error {
	if err := r.db.WithContext(ctx).Create(inscription).Error; err != nil {
		r.logger.Error("[INSCRIPTION-API] Error al crear inscripción en la base de datos", zap.Error(err))
		return errors.ErrInternalServer
	}
	return nil
}

func (r *gormInscriptionRepository) GetMyCourses(ctx context.Context, userId string) ([]string, error) {
	var inscriptions []models.Inscripto
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&inscriptions).Error; err != nil {
		r.logger.Error("[INSCRIPTION-API] Error al obtener cursos del usuario", zap.String("user_id", userId), zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	if len(inscriptions) == 0 {
		return nil, nil
	}

	courseIds := make([]string, len(inscriptions))
	for i, inscription := range inscriptions {
		courseIds[i] = inscription.CourseId
	}
	return courseIds, nil
}

func (r *gormInscriptionRepository) GetStudentsInCourse(ctx context.Context, courseId string) ([]string, error) {
	var inscriptions []models.Inscripto
	if err := r.db.WithContext(ctx).Where("course_id = ?", courseId).Find(&inscriptions).Error; err != nil {
		r.logger.Error("[INSCRIPTION-API] Error al obtener estudiantes del curso", zap.String("course_id", courseId), zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	if len(inscriptions) == 0 {
		return nil, nil
	}

	studentIds := make([]string, len(inscriptions))
	for i, inscription := range inscriptions {
		studentIds[i] = inscription.UserId
	}
	return studentIds, nil
}

func (r *gormInscriptionRepository) IsEnrolled(ctx context.Context, courseId string, userId string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Inscripto{}).
		Where("course_id = ? AND user_id = ?", courseId, userId).
		Count(&count).Error; err != nil {
		r.logger.Error("[INSCRIPTION-API] Error al verificar si el usuario está inscrito", zap.String("user_id", userId), zap.String("course_id", courseId), zap.Error(err))
		return false, errors.ErrInternalServer
	}
	enrolled := count > 0
	return enrolled, nil
}
