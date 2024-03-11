package repository

import (
	"context"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
)

type Managers interface {
	Create(ctx context.Context, manager models.ManagerBase) (int, error)
	GetByLogin(ctx context.Context, managerLogin string) (models.Manager, error)
}

type Specialists interface {
	Create(ctx context.Context, specialist models.SpecialistCreate) (int, error)
	GetByID(ctx context.Context, specialistID int) (models.Specialist, error)
	GetByLogin(ctx context.Context, specialistLogin string) (models.Specialist, error)
	Update(ctx context.Context, specialistUpdate models.Specialist, newPasswordFlag bool) error
	Delete(ctx context.Context, specialistID int) error
}

type Cameras interface {
	Create(ctx context.Context, camera models.CameraBase) (string, error)
	Get(ctx context.Context, cameraID string) (models.Camera, error)
	Delete(ctx context.Context, cameraID string) error
}

type Cases interface {
	CreateCase(ctx context.Context, caseData models.CaseBase) (int, error)
	GetCaseByID(ctx context.Context, caseID int) (models.Case, error)
	GetCasesByLevel(ctx context.Context, specialistID, level, cursor int) (models.CaseCursor, error)
	DeleteCase(ctx context.Context, caseID int) error

	CreateRated(ctx context.Context, rated models.RatedBase) (int, error)
	GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error)
	UpdateRatedStatus(ctx context.Context, newRated models.RatedUpdate) error

	GetFulCaseByID(ctx context.Context, caseID int) (models.CaseFul, error)
}

type Violations interface {
	Create(violations []models.Violation) (int, error)
}

type Contacts interface {
	Create(contacts []models.Contact) (int, error)
}
