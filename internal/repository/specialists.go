package repository

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type specialistsRepo struct {
	db *sqlx.DB
}

func InitSpecialistsRepo(db *sqlx.DB) Specialists {
	return specialistsRepo{
		db: db,
	}
}

func (s specialistsRepo) Create(ctx context.Context, specialist models.SpecialistCreate) (int, error) {
	var createdSpecialistID int

	tx, err := s.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	hashedPassword := utils.HashPassword(specialist.Password)

	specialistCreateQuery := `INSERT INTO specialists (login, hashed_password, fullname, photo_url)
						VALUES ($1, $2, $3, $4)
						RETURNING id;`

	err = tx.QueryRowxContext(ctx, specialistCreateQuery,
		specialist.Login, hashedPassword, specialist.Fullname, specialist.PhotoUrl).Scan(&createdSpecialistID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdSpecialistID, nil
}

func (s specialistsRepo) Get(ctx context.Context, specialistID int) (models.Specialist, error) {
	var specialist models.Specialist

	specialistGetQuery := `SELECT id, login, hashed_password, fullname, level, photo_url, is_verified
						FROM specialists
						WHERE id=$1;`

	err := s.db.GetContext(ctx, &specialist, specialistGetQuery, specialistID)
	if err != nil {
		return models.Specialist{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	return specialist, nil
}

func (s specialistsRepo) Update(ctx context.Context, specialistUpdate models.SpecialistUpdate) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	specialistUpdate.Password = string(utils.HashPassword(specialistUpdate.Password))

	specialistUpdateQuery := `UPDATE specialists
							  SET login = :login, hashed_password = :hashed_password, fullname = :fullname, level = :level, photo_url = :photo_url, is_verified = :is_verified
							  WHERE id = :id;`

	res, err := tx.NamedExecContext(ctx, specialistUpdateQuery, specialistUpdate)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}

func (s specialistsRepo) Delete(ctx context.Context, specialistID int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	specialistDeleteQuery := `DELETE FROM specialists WHERE id=$1;`

	res, err := tx.ExecContext(ctx, specialistDeleteQuery, specialistID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}
