package exfuncs

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func LoadViolation() error {
	var violations []models.Violation

	config.InitConfig()

	db := database.GetDB()
	violationRepo := repository.InitViolationRepo(db)

	f, err := excelize.OpenFile("../internal/exloads/exfiles/violations_example.xlsx")
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		var violation models.Violation

		amount, err := strconv.Atoi(row[1])
		if err != nil {
			return fmt.Errorf("Ошибка при преобразовании строки в число в строке %d: %v\n", i+1, err)
		}
		violation.Type = row[0]
		violation.Amount = amount

		violations = append(violations, violation)
	}

	num, err := violationRepo.Create(violations)
	if err != nil {
		return err
	}

	fmt.Printf("Вставлено %d строк.", num)

	return nil
}
