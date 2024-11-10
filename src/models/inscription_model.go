// models/inscripto.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inscripto struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CourseId  string         `gorm:"type:string;not null" json:"course_id"`
	UserId    string         `gorm:"type:string;not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
