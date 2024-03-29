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
	managerRepo    repository.Managers
	specialistRepo repository.Specialists
	cameraRepo     repository.Cameras
	caseRepo       repository.Cases
	dbResponseTime time.Duration
	logger         *log.Logs
}

func InitPublicService(
	managerRepo repository.Managers,
	specialistRepo repository.Specialists,
	cameraRepo repository.Cameras,
	caseRepo repository.Cases,
	logger *log.Logs,
) Public {
	return publicService{
		managerRepo:    managerRepo,
		specialistRepo: specialistRepo,
		cameraRepo:     cameraRepo,
		caseRepo:       caseRepo,
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

	p.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessCreate, "specialist", createdSpecialistID))
	return createdSpecialistID, nil
}

func (p publicService) ManagerLogin(ctx context.Context, manager models.ManagerBase) (bool, models.Manager, error) {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	managerData, err := p.managerRepo.GetByLogin(ctx, manager.Login)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return false, models.Manager{}, err
	}

	isCompare := utils.ComparePassword(managerData.Password, manager.Password)

	return isCompare, managerData, nil
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

func (p publicService) CameraCreate(ctx context.Context, camera models.CameraBase) (string, error) {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	createdCameraID, err := p.cameraRepo.Create(ctx, camera)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return "", err
	}

	p.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessCreate, "camera", createdCameraID))
	return createdCameraID, nil
}

func (p publicService) CameraDelete(ctx context.Context, cameraID string) error {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	err := p.cameraRepo.Delete(ctx, cameraID)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}

	p.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessDelete, "camera", cameraID))
	return nil
}

func (p publicService) CaseCreate(ctx context.Context, caseData models.CaseBase) (int, error) {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	createdCaseID, err := p.caseRepo.CreateCase(ctx, caseData)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return 0, err
	}

	p.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessCreate, "case", createdCaseID))
	return createdCaseID, nil
}

func (p publicService) CaseDelete(ctx context.Context, caseID int) error {
	ctx, cansel := context.WithTimeout(ctx, p.dbResponseTime)
	defer cansel()

	err := p.caseRepo.DeleteCase(ctx, caseID)
	if err != nil {
		p.logger.ErrorLogger.Error().Msg(err.Error())
		return err
	}

	p.logger.InfoLogger.Info().Msg(fmt.Sprintf(responses.ResponseSuccessDelete, "camera", caseID))
	return nil
}
