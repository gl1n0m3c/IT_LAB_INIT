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
