package services

import (
	"context"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
)

type Managers interface {
	GetFulCaseByID(ctx context.Context, caseID int) (models.CaseFul, error)
}

type Public interface {
	ManagerLogin(ctx context.Context, manager models.ManagerBase) (bool, models.Manager, error)

	SpecialistRegister(ctx context.Context, specialist models.SpecialistCreate) (int, error)
	SpecialistLogin(ctx context.Context, specialist models.SpecialistLogin) (bool, models.Specialist, error)

	CameraCreate(ctx context.Context, camera models.CameraBase) (string, error)
	CameraDelete(ctx context.Context, cameraID string) error

	CaseCreate(ctx context.Context, caseData models.CaseBase) (int, error)
	CaseDelete(ctx context.Context, caseID int) error
}

type Specialists interface {
	GetMe(ctx context.Context, specialistID int) (models.Specialist, error)
	UpdateMe(ctx context.Context, specialistUpdate models.SpecialistUpdate) error

	GetCasesByLevel(ctx context.Context, specialistID, cursor int) (models.CaseCursor, error)

	CreateRated(ctx context.Context, rated models.RatedBase) (int, error)
	GetRatedSolved(ctx context.Context, specialistID, cursor int) (models.RatedCursor, error)
	UpdateRatedStatus(ctx context.Context, specialistID int, newRated models.RatedUpdate) error
}
