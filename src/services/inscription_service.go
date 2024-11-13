package services

import (
	"context"
	"inscription-api/src/clients"
	"inscription-api/src/dto"
	"inscription-api/src/errors"
	"inscription-api/src/models"

	"go.uber.org/zap"
)

type InscriptionService interface {
	EnrollStudent(ctx context.Context, courseId string, userId string) (dto.EnrollRequestResponseDto, error)
	GetMyCourses(ctx context.Context, userId string) (dto.MyCourses, error)
	GetStudentsInCourse(ctx context.Context, courseId string) (dto.StudentsInCourse, error)
	IsEnrolled(ctx context.Context, courseId string, userId string) (bool, error)
}

type inscriptionService struct {
	repo   clients.InscriptionRepository
	logger *zap.Logger
}

func NewInscriptionService(repo clients.InscriptionRepository, logger *zap.Logger) InscriptionService {
	return &inscriptionService{repo: repo, logger: logger}
}

func (s *inscriptionService) EnrollStudent(ctx context.Context, courseId string, userId string) (dto.EnrollRequestResponseDto, error) {
	enrolled, err := s.IsEnrolled(ctx, courseId, userId)
	if err != nil {
		return dto.EnrollRequestResponseDto{}, err
	}
	if enrolled {
		s.logger.Warn("[INSCRIPTION-API] El estudiante ya está inscrito", zap.String("user_id", userId), zap.String("course_id", courseId))
		return dto.EnrollRequestResponseDto{}, errors.ErrDuplicateEnroll
	}

	inscription := &models.Inscripto{
		CourseId: courseId,
		UserId:   userId,
	}
	Inscription, err := s.repo.Create(ctx, inscription)
	if err != nil {
		return dto.EnrollRequestResponseDto{}, err
	}
	return dto.EnrollRequestResponseDto{
		CourseId: Inscription.CourseId,
		UserId:   Inscription.UserId,
	}, nil
}

func (s *inscriptionService) GetMyCourses(ctx context.Context, userId string) (dto.MyCourses, error) {
	if userId == "" {
		s.logger.Warn("[INSCRIPTION-API] El ID de usuario es requerido")
		return nil, errors.ErrMissingUserId
	}

	courseIds, err := s.repo.GetMyCourses(ctx, userId)
	if err != nil {
		return nil, err
	}

	if courseIds == nil {
		return nil, errors.ErrNoResults
	}

	courses := make(dto.MyCourses, len(courseIds))
	for i, courseId := range courseIds {
		courses[i] = dto.MyCourse{
			CourseId: courseId,
		}
	}
	return courses, nil
}

func (s *inscriptionService) GetStudentsInCourse(ctx context.Context, courseId string) (dto.StudentsInCourse, error) {
	studentIds, err := s.repo.GetStudentsInCourse(ctx, courseId)
	if err != nil {
		s.logger.Error("[INSCRIPTION-API] Error al obtener los estudiantes del curso", zap.Error(err))
		return nil, err
	}

	if studentIds == nil {
		return nil, errors.ErrNoResults
	}

	// Aquí podrías hacer una llamada al servicio de usuarios para obtener más detalles
	students := make(dto.StudentsInCourse, len(studentIds))
	for i, userId := range studentIds {
		students[i] = dto.Student{
			UserId: userId,
			// Agrega otros campos si tienes más información
		}
	}

	return students, nil
}

func (s *inscriptionService) IsEnrolled(ctx context.Context, courseId string, userId string) (bool, error) {
	enrolled, err := s.repo.IsEnrolled(ctx, courseId, userId)
	if err != nil {
		s.logger.Error("[INSCRIPTION-API] Error al verificar inscripción", zap.Error(err))
		return false, err
	}
	return enrolled, nil
}
