package services

import (
	"context"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
)

type Public interface {
	SpecialistRegister(ctx context.Context, specialist models.SpecialistCreate) (int, error)
	SpecialistLogin(ctx context.Context, specialist models.SpecialistLogin) (bool, models.Specialist, error)

	CameraCreate(ctx context.Context, camera models.CameraBase) (int, error)
	CameraDelete(ctx context.Context, cameraID int) error
}
