package repository

import (
	"context"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"time"
)

type Managers interface {
	Create(ctx context.Context, manager models.ManagerBase) (int, error)
	GetByLogin(ctx context.Context, managerLogin string) (models.Manager, error)
}

type Specialists interface {
	Create(ctx context.Context, specialist models.SpecialistCreate) (int, error)
	GetByID(ctx context.Context, specialistID int) (models.Specialist, error)
	GetByLogin(ctx context.Context, specialistLogin string) (models.Specialist, error)
	GetSpecialistRating(ctx context.Context, timeStart, timeEnd time.Time, cursor int) (models.RatingSpecialistCountCursor, error)
	GetFulRating(ctx context.Context) ([]models.RatingSpecialistFul, error)
	GetOnlyRating(timeStart, timeEnd time.Time) ([]models.RatingSpecialistID, error)
	UpdateSpecialistsIncDecLevel(incrementIDs, decrementIDs []int) error
	UpdateMain(ctx context.Context, specialistUpdate models.Specialist, newPasswordFlag bool) error
	Delete(ctx context.Context, specialistID int) error
}

type Cameras interface {
	Create(ctx context.Context, camera models.CameraBase) (string, error)
	Get(ctx context.Context, cameraID string) (models.Camera, error)
	Delete(ctx context.Context, cameraID string) error
}

type Cases interface {
	CreateCase(ctx context.Context, caseData models.CaseBase) (int, error)
	UpdateCaseLevel(ctx context.Context, caseID, level int) error
	UpdateCaseSetSolved(ctx context.Context, caseID int, rightChoice bool) error
	GetFineData(ctx context.Context, caseID int) (models.FineData, error)
	GetCaseLevelSolvedRatingsTrueByID(ctx context.Context, caseID int) (int, int, int, bool, error)
	GetCasesByLevel(ctx context.Context, specialistID, level, cursor int) (models.CaseCursor, error)
	DeleteCase(ctx context.Context, caseID int) error

	CreateRated(ctx context.Context, rated models.RatedBase) (int, error)
	GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error)
	GetNumberRatedByCaseID(ctx context.Context, caseID int) (int, error)

	GetFulCaseByID(ctx context.Context, caseID int) (models.CaseFul, error)
}

type Violations interface {
	Create(violations []models.Violation) (int, error)
}

type Contacts interface {
	Create(contacts []models.Contact) (int, error)
}
