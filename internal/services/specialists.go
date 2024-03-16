package services

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/sender"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/guregu/null"
	"github.com/spf13/viper"
	"os"
	"time"
)

type specialistService struct {
	specialistRepo repository.Specialists
	caseRepo       repository.Cases
	k              int
	dbResponseTime time.Duration
	logger         *log.Logs
}

func InitSpecialistService(
	specialistRepo repository.Specialists,
	caseRepo repository.Cases,
	logger *log.Logs,
) Specialists {
	return specialistService{
		specialistRepo: specialistRepo,
		caseRepo:       caseRepo,
		k:              viper.GetInt(config.K),
		dbResponseTime: time.Duration(viper.GetInt(config.DBResponseTime)) * time.Second,
		logger:         logger,
	}
}

func (s specialistService) GetMe(ctx context.Context, specialistID int) (models.Specialist, error) {
	ctx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	specialist, err := s.specialistRepo.GetByID(ctx, specialistID)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return models.Specialist{}, err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessGet, "specialist"))

	return specialist, nil
}

func (s specialistService) UpdateMe(ctx context.Context, specialistUpdate models.SpecialistUpdate) error {
	var passwordFlag bool

	getCtx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	specialist, err := s.specialistRepo.GetByID(getCtx, specialistUpdate.ID)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}

	if specialistUpdate.FullName != "" {
		specialist.Fullname = null.NewString(specialistUpdate.FullName, true)
	}

	if !utils.ComparePassword(specialist.Password, specialistUpdate.Password) {
		specialist.Password = specialistUpdate.Password
		passwordFlag = true
	}

	if specialistUpdate.PhotoUrl != "" {
		if specialist.PhotoUrl.Valid {
			err = os.Remove("../" + specialist.PhotoUrl.String)
			if err != nil {
				return err
			}
		}

		specialist.PhotoUrl = null.NewString(specialistUpdate.PhotoUrl, true)
	}

	updCtx, updCansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer updCansel()

	err = s.specialistRepo.UpdateMain(updCtx, specialist, passwordFlag)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessUpdate, "specialist"))

	return nil
}

func (s specialistService) CreateRated(ctx context.Context, rated models.RatedBase) (int, error) {
	specCtx, specCansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer specCansel()

	// Проверка, что аккаунт специалиста подтвержден
	specialist, err := s.specialistRepo.GetByID(specCtx, rated.SpecialistID)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return 0, err
	}
	if !specialist.IsVerified {
		s.logger.ErrorLogger.Info().Msg(customErrors.UserUnverified.Error())
		return 0, customErrors.UserUnverified
	}

	caseCtx1, caseCansel1 := context.WithTimeout(ctx, s.dbResponseTime)
	defer caseCansel1()

	// Проверка, что уровень специалиста совпадает с уровнем кейса + кейс еще не разрешен + проверка на количество оценок
	level, numberOfRated, numberOfTrue, isSolved, err := s.caseRepo.GetCaseLevelSolvedRatingsTrueByID(caseCtx1, rated.CaseID, specialist.Level)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return 0, err
	}
	if isSolved || numberOfRated >= s.k {
		s.logger.ErrorLogger.Info().Msg(customErrors.CaseAlreadySolved.Error())
		return 0, customErrors.CaseAlreadySolved
	}

	caseCtx2, caseCansel2 := context.WithTimeout(ctx, s.dbResponseTime)
	defer caseCansel2()

	createdRatedID, err := s.caseRepo.CreateRated(caseCtx2, rated)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return 0, err
	}

	// Проверка на консенсус
	if s.k-1 == numberOfRated {
		if (numberOfRated == numberOfTrue && rated.Choice) || (numberOfTrue == 0 && !rated.Choice) {
			var rightChoice bool
			if numberOfRated == numberOfTrue {
				rightChoice = true
			}

			updateCtx, updateCansel := context.WithTimeout(ctx, s.dbResponseTime)
			defer updateCansel()

			err := s.caseRepo.UpdateCaseSetSolved(updateCtx, rated.CaseID, rightChoice)
			if err != nil {
				s.logger.ErrorLogger.Error().Msg(err.Error())
				return 0, err
			}

			fineDataCtx, fineDataCansel := context.WithTimeout(ctx, s.dbResponseTime)
			defer fineDataCansel()

			fineData, err := s.caseRepo.GetFineData(fineDataCtx, rated.CaseID)
			if err != nil {
				s.logger.ErrorLogger.Error().Msg(err.Error())
				return 0, err
			}

			err = sender.MailSender(fineData)
			if err != nil {
				s.logger.ErrorLogger.Error().Msg(err.Error())
				return 0, err
			}
		} else {
			updateCtx, updateCansel := context.WithTimeout(ctx, s.dbResponseTime)
			defer updateCansel()

			err := s.caseRepo.UpdateCaseLevel(updateCtx, rated.CaseID, level+1)
			if err != nil {
				s.logger.ErrorLogger.Error().Msg(err.Error())
				return 0, err
			}
		}
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessCreate, "rated_case", createdRatedID))

	return createdRatedID, nil
}

func (s specialistService) GetCasesByLevel(ctx context.Context, specialistID, cursor int) (models.CaseCursor, error) {
	specCtx, specCansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer specCansel()

	// Проверка, что аккаунт специалиста подтвержден
	specialist, err := s.specialistRepo.GetByID(specCtx, specialistID)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return models.CaseCursor{}, err
	}
	if !specialist.IsVerified {
		s.logger.ErrorLogger.Info().Msg(customErrors.UserUnverified.Error())
		return models.CaseCursor{}, customErrors.UserUnverified
	}

	ctx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	cases, err := s.caseRepo.GetCasesByLevel(ctx, specialist.ID, specialist.Level, cursor)

	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return models.CaseCursor{}, err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessGet, "cases_by_level"))

	return cases, nil
}

func (s specialistService) GetRatedSolved(ctx context.Context, specialistID, cursor int) (models.RatedCursor, error) {
	specCtx, specCansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer specCansel()

	// Проверка, что аккаунт специалиста подтвержден
	specialist, err := s.specialistRepo.GetByID(specCtx, specialistID)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return models.RatedCursor{}, err
	}
	if !specialist.IsVerified {
		s.logger.ErrorLogger.Info().Msg(customErrors.UserUnverified.Error())
		return models.RatedCursor{}, customErrors.UserUnverified
	}

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

func (s specialistService) UpdateRatedStatus(ctx context.Context, specialistID int, newRated models.RatedUpdate) error {
	specCtx, specCansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer specCansel()

	// Проверка, что аккаунт специалиста подтвержден
	specialist, err := s.specialistRepo.GetByID(specCtx, specialistID)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}
	if !specialist.IsVerified {
		s.logger.ErrorLogger.Info().Msg(customErrors.UserUnverified.Error())
		return customErrors.UserUnverified
	}

	ctx, cansel := context.WithTimeout(ctx, s.dbResponseTime)
	defer cansel()

	err = s.caseRepo.UpdateRatedStatus(ctx, newRated)
	if err != nil {
		s.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}

	s.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessUpdate, "rated_case"))

	return nil

}
