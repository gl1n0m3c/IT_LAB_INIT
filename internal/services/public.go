package services

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/spf13/viper"
	"time"
)

type publicService struct {
	specialistRepo repository.Specialists
	dbResponseTime time.Duration
	logger         *log.Logs
}

func InitPublicService(
	specialistRepo repository.Specialists,
	logger *log.Logs,
) Public {
	return publicService{
		specialistRepo: specialistRepo,
		dbResponseTime: time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second,
		logger:         logger,
	}
}

func (p publicService) SpecialistRegister(ctx context.Context, specialist models.SpecialistCreate) (int, error) {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	createdSpecialistID, err := p.specialistRepo.Create(ctx, specialist)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return 0, err
	}

	p.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.Response201, "specialist", createdSpecialistID))
	return createdSpecialistID, nil
}

func (p publicService) SpecialistLogin(ctx context.Context, specialist models.SpecialistLogin) (bool, models.Specialist, error) {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	specialistData, err := p.specialistRepo.GetByLogin(ctx, specialist.Login)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return false, models.Specialist{}, err
	}

	isCompare := utils.ComparePassword(specialistData.Password, specialist.Password)

	return isCompare, specialistData, nil
}
