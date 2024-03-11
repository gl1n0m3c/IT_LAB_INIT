package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type managerRepo struct {
	db *sqlx.DB
}

func InitManagerRepo(db *sqlx.DB) Managers {
	return managerRepo{db: db}
}

func (m managerRepo) Create(ctx context.Context, manager models.ManagerBase) (int, error) {
	var createdManagerID int

	tx, err := m.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	hashedPassword := utils.HashPassword(manager.Password)

	specialistCreateQuery := `INSERT INTO managers (login, hashed_password)
						VALUES ($1, $2)
						RETURNING id;`

	err = tx.QueryRowxContext(ctx, specialistCreateQuery, manager.Login, hashedPassword).Scan(&createdManagerID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, customErrors.UniqueSpecialistErr
		}

		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdManagerID, nil
}

func (m managerRepo) GetByLogin(ctx context.Context, managerLogin string) (models.Manager, error) {
	var manager models.Manager

	managerGetQuery := `SELECT id, login, hashed_password
						FROM managers
						WHERE login=$1;`

	err := m.db.GetContext(ctx, &manager, managerGetQuery, managerLogin)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Manager{}, customErrors.NoRowsSpecialistLoginErr
		default:
			return models.Manager{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
	}

	return manager, nil
}
