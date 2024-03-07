package services

import (
	"context"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
)

type Public interface {
	SpecialistRegister(ctx context.Context, specialist models.SpecialistCreate) (int, error)
	SpecialistLogin(ctx context.Context, specialist models.SpecialistLogin) (bool, models.Specialist, error)

	CameraCreate(ctx context.Context, camera models.CameraBase) (string, error)
	CameraDelete(ctx context.Context, cameraID string) error

	CaseCreate(ctx context.Context, caseData models.CaseBase) (int, error)
	CaseDelete(ctx context.Context, caseID int) error
}

type Specialists interface {
	CreateRated(ctx context.Context, rated models.RatedBase) (int, error)
	GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error)
	UpdateRatedStatus(ctx context.Context, newRated models.RatedUpdate) error
}
