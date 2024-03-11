package services

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/spf13/viper"
	"time"
)

type managerService struct {
	caseRepo        repository.Cases
	specialistsRepo repository.Specialists
	dbResponseTime  time.Duration
	logger          *log.Logs
}

func InitManagerService(
	caseRepo repository.Cases,
	specialistsRepo repository.Specialists,
	logger *log.Logs,
) Managers {
	return managerService{
		caseRepo:        caseRepo,
		specialistsRepo: specialistsRepo,
		dbResponseTime:  time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second,
		logger:          logger,
	}
}

func (m managerService) GetFulCaseByID(ctx context.Context, caseID int) (models.CaseFul, error) {
	ctx, cansel := context.WithTimeout(ctx, m.dbResponseTime)
	defer cansel()

	caseFulData, err := m.caseRepo.GetFulCaseByID(ctx, caseID)
	if err != nil {
		m.logger.ErrorLogger.Error().Msg(err.Error())
		return models.CaseFul{}, err
	}

	m.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessGet, "case_ful"))

	return caseFulData, nil
}

func (m managerService) GetSpecialistRating(ctx context.Context, timeStart, timeEnd time.Time, cursor int) (models.RatingSpecialistCountCursor, error) {
	ctx, cansel := context.WithTimeout(ctx, m.dbResponseTime)
	defer cansel()

	specialists, err := m.specialistsRepo.GetSpecialistRating(ctx, timeStart, timeEnd, cursor)
	if err != nil {
		m.logger.ErrorLogger.Error().Msg(err.Error())
		return models.RatingSpecialistCountCursor{}, err
	}

	m.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessGet, "specialists_raiting"))

	return specialists, nil
}
