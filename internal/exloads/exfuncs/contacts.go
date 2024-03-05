package exfuncs

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/xuri/excelize/v2"
)

func LoadContacts() error {
	var contacts []models.Contact
	var keys []string

	config.InitConfig()

	db := database.GetDB()
	contactRepo := repository.InitContactRepo(db)

	f, err := excelize.OpenFile("../internal/exloads/exfiles/contacts_example.xlsx")
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
			for j, key := range row {
				if j == 0 {
					continue
				}
				keys = append(keys, key)
			}
			continue
		}

		var contact models.Contact
		userContacts := make(map[string]string)

		for j, key := range row {
			if j == 0 {
				continue
			}
			userContacts[keys[j-1]] = key
		}

		contact.Transport = row[0]
		contact.UserContacts = userContacts

		contacts = append(contacts, contact)
	}

	num, err := contactRepo.Create(contacts)
	if err != nil {
		return err
	}

	fmt.Printf("Вставлено %d строк.", num)

	return nil
}
