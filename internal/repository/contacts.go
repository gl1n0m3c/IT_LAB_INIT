package repository

import (
	"encoding/json"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type contactRepo struct {
	db *sqlx.DB
}

func InitContactRepo(db *sqlx.DB) Contacts {
	return contactRepo{db: db}
}

func (c contactRepo) Create(contacts []models.Contact) (int, error) {
	var accepted int

	tx, err := c.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	contactsQueue := `INSERT INTO contacts (transport, contacts)
					  VALUES ($1, $2)`

	for i, contact := range contacts {
		savepoint := fmt.Sprintf("sp_%d", i)
		if _, err := tx.Exec("SAVEPOINT " + savepoint); err != nil {
			return 0, err
		}

		userContactsJSON, err := json.Marshal(contact.UserContacts)
		if err != nil {
			return 0, err
		}

		res, err := tx.Exec(contactsQueue, contact.Transport, userContactsJSON)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				_, rbErr := tx.Exec("ROLLBACK TO SAVEPOINT " + savepoint)
				if rbErr != nil {
					return 0, rbErr
				}
				continue
			}

			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, fmt.Errorf("error: %v", err.Error()+rbErr.Error())
			}
			return 0, err
		}

		count, err := res.RowsAffected()
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, fmt.Errorf("error: %v", err.Error()+rbErr.Error())
			}
			return 0, err
		}

		if count == 0 {
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, fmt.Errorf("error: %v", err.Error()+rbErr.Error())
			}
			return 0, fmt.Errorf("no rows were inserted")
		}

		accepted++
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return accepted, nil
}
