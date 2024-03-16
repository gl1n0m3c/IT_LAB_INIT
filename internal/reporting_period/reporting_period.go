package reporting_period

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"time"
)

func StartReporting(db *sqlx.DB, logger *log.Logs) {
	// Запуск горутины, обновляющей уровень специалистов по истечении отчетного периода
	// Закоменченная шняга для демонстрации
	go func() {
		specialistsRepo := repository.InitSpecialistsRepo(db)
		//timeStart := time.Now().Add(time.Hour * 24 * 360 * -2)
		timeStart := time.Now()
		timeEnd := time.Now().Add(time.Duration(viper.GetInt(config.ReportingPeriod)) * time.Hour * 24)
		ticker := time.NewTicker(time.Duration(viper.GetInt(config.ReportingPeriod)) * time.Hour * 24)
		//timeEnd := time.Now().Add(time.Second * 10)
		//ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				var specialistIncrementIDs []int
				var specialistDecrementIDs []int

				rating, err := specialistsRepo.GetRating(timeStart, timeEnd)
				if err != nil {
					logger.ErrorLogger.Error().Msg(fmt.Sprintf("Ошибка при обновлении уровней специалистов:", err))
				}

				tenProc := (len(rating) + 10 - 1) / 10

				for i := 0; i < tenProc; i++ {
					specialistIncrementIDs = append(specialistIncrementIDs, rating[i].ID)
				}
				var skipped int
				for i := 0; i < tenProc+skipped; i++ {
					if rating[len(rating)-i-1].Level == 1 {
						skipped++
						continue
					}

					var flag bool
					for _, val := range specialistIncrementIDs {
						if val == rating[len(rating)-i-1].ID {
							flag = true
							break
						}
					}
					if flag {
						break
					}

					specialistDecrementIDs = append(specialistDecrementIDs, rating[len(rating)-i-1].ID)
				}

				err = specialistsRepo.UpdateSpecialistsIncDecLevel(specialistIncrementIDs, specialistDecrementIDs)
				if err != nil {
					logger.ErrorLogger.Error().Msg(fmt.Sprintf("Ошибка при обновлении уровней специалистов:", err))
				}
				logger.ErrorLogger.Info().Msg("Уровни специалистов обновлены")

				timeStart.Add(time.Duration(viper.GetInt(config.ReportingPeriod)) * time.Hour * 24)
				timeEnd.Add(time.Duration(viper.GetInt(config.ReportingPeriod)) * time.Hour * 24)
			}
		}
	}()

}
