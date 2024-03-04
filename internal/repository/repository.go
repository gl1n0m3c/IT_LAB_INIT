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
	Create(ctx context.Context, camera models.CameraBase) (int, error)
	Get(ctx context.Context, cameraID int) (models.Camera, error)
	Delete(ctx context.Context, cameraID int) error
}

type Cases interface {
	Create(ctx context.Context, caseData models.CaseBase) (int, error)
	Delete(ctx context.Context, caseID int) error
}
