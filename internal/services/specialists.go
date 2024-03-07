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

type specialistService struct {
	caseRepo       repository.Cases
	dbResponseTime time.Duration
	logger         *log.Logs
}

func InitSpecialistService(
	caseRepo repository.Cases,
	logger *log.Logs,
) Specialists {
	return specialistService{
		caseRepo:       caseRepo,
		dbResponseTime: time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second,
		logger:         logger,
	}
}

func (s specialistService) CreateRated(ctx context.Context, rated models.RatedBase) (int, error) {
	ctx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	createdRatedID, err := s.caseRepo.CreateRated(ctx, rated)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return 0, err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessCreate, "rated_case", createdRatedID))

	return createdRatedID, nil

}

func (s specialistService) GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error) {
	ctx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	solvedRated, err := s.caseRepo.GetRatedSolved(ctx, cursor)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return models.RatedCursor{}, err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessGet, "rated_cases"))

	return solvedRated, nil

}

func (s specialistService) UpdateRatedStatus(ctx context.Context, newRated models.RatedUpdate) error {
	ctx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	err := s.caseRepo.UpdateRatedStatus(ctx, newRated)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessUpdate, "rated_case"))

	return nil

}
