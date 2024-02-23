package specialists

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
)

type specialistsRepo struct {
	db *sql.DB
}

func InitSpecialistsRepo(db *sql.DB) repository.Specialists {
	return specialistsRepo{
		db: db,
	}
}

func (s specialistsRepo) Create(ctx context.Context, specialist models.SpecialistCreate) (int, error) {
	var createdSpecialistID int

	tx, err := s.db.Begin()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	hashedPassword := utils.HashPassword(specialist.Password)

	specialistCreateQuery := `INSERT INTO specialists (login, hashed_password, fullname, photo_url)
						VALUES ($1, $2, $3, $4)
						RETURNING id;`

	row := tx.QueryRowContext(ctx, specialistCreateQuery,
		specialist.Login, hashedPassword, specialist.Fullname, specialist.PhotoUrl)

	if err = row.Scan(&createdSpecialistID); err != nil {
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

	specialistGetQuery := `SELECT (id, login, hashed_password, fullname, level, photo_url, is_verified)
						FROM specialists
						WHERE id=$1;`

	row := s.db.QueryRowContext(ctx, specialistGetQuery, specialistID)

	// В поле password кладется хешированынй пароль
	err := row.Scan(specialist.ID, specialist.Login, specialist.Password, specialist.Fullname,
		specialist.Level, specialist.PhotoUrl, specialist.IsVerified)
	if err != nil {
		return models.Specialist{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	return specialist, nil
}

func (s specialistsRepo) Update(ctx context.Context, specialistUpdate models.SpecialistUpdate) error {
	tx, err := s.db.Begin()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	specialistUpdateQuery := `UPDATE specialists
							  SET login = $1, hashed_password = $2, fullname = $3, level = $4, photo_url = $5, is_verified = $6
							  WHERE id = $7;`

	hashedPassword := utils.HashPassword(specialistUpdate.Password)

	res, err := tx.ExecContext(ctx, specialistUpdateQuery,
		specialistUpdate.Login, hashedPassword, specialistUpdate.Fullname,
		specialistUpdate.Level, specialistUpdate.PhotoUrl, specialistUpdate.IsVerified,
		specialistUpdate.ID)
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
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf("too few lines returned `%v`", count)})
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}

func (s specialistsRepo) Delete(ctx context.Context, specialistID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	specialistDeleteQuery := `DELETE FROM specialists where id=$1`

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
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf("too few lines returned `%v`", count)})
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}
