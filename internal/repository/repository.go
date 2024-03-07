package repository

import (
	"context"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
)

type Specialists interface {
	Create(ctx context.Context, specialist models.SpecialistCreate) (int, error)
	GetByID(ctx context.Context, specialistID int) (models.Specialist, error)
	GetByLogin(ctx context.Context, specialistLogin string) (models.Specialist, error)
	Update(ctx context.Context, specialistUpdate models.SpecialistUpdate) error
	Delete(ctx context.Context, specialistID int) error
}

type Cameras interface {
	Create(ctx context.Context, camera models.CameraBase) (string, error)
	Get(ctx context.Context, cameraID string) (models.Camera, error)
	Delete(ctx context.Context, cameraID string) error
}

type Cases interface {
	CreateCase(ctx context.Context, caseData models.CaseBase) (int, error)
	GetCasesByLevel(ctx context.Context, level, cursor int) (models.CaseCursor, error)
	DeleteCase(ctx context.Context, caseID int) error

	CreateRated(ctx context.Context, rated models.RatedBase) (int, error)
	GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error)
	UpdateRatedStatus(ctx context.Context, newRated models.RatedUpdate) error
}

type Violations interface {
	Create(violations []models.Violation) (int, error)
}

type Contacts interface {
	Create(contacts []models.Contact) (int, error)
}
