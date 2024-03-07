package repository

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type violationRepo struct {
	db *sqlx.DB
}

func InitViolationRepo(db *sqlx.DB) Violations {
	return violationRepo{db: db}
}

func (v violationRepo) Create(violations []models.Violation) (int, error) {
	var accepted int

	tx, err := v.db.Beginx()
	if err != nil {
		return 0, err
	}

	cameraQueue := `INSERT INTO violations (id, type, amount)
					VALUES ($1, $2, $3)`

	for i, violation := range violations {
		savepoint := fmt.Sprintf("sp_%d", i)
		if _, err := tx.Exec("SAVEPOINT " + savepoint); err != nil {
			return 0, err
		}

		uuidBytes, err := uuid.NewV4()
		if err != nil {
			return 0, err
		}

		key := uuidBytes.String()

		res, err := tx.Exec(cameraQueue, key, violation.Type, violation.Amount)
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
