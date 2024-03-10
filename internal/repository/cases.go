package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/custom_errors"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type caseRepo struct {
	db              *sqlx.DB
	CasesPerRequest int
}

func InitCaseRepo(
	db *sqlx.DB,
	CasesPerRequest int,
) Cases {
	return caseRepo{
		db:              db,
		CasesPerRequest: CasesPerRequest,
	}
}

func (c caseRepo) CreateCase(ctx context.Context, caseData models.CaseBase) (int, error) {
	var createdCaseID int

	tx, err := c.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	caseCreateQuery := `INSERT INTO cases (camera_id, transport, violation_id, violation_value, level, datetime, photo_url)
						VALUES ($1, $2, $3, $4, $5, $6, $7)
						RETURNING id;`

	err = tx.QueryRowxContext(ctx, caseCreateQuery,
		caseData.CameraID, caseData.Transport, caseData.ViolationID, caseData.ViolationValue,
		caseData.Level, caseData.Datetime, caseData.PhotoUrl).Scan(&createdCaseID)
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

	return createdCaseID, nil
}

func (c caseRepo) GetCaseByID(ctx context.Context, caseID int) (models.Case, error) {
	var caseData models.Case

	caseGetQuery := `SELECT id, camera_id, transport, violation_id,violation_value, level, datetime, photo_url
					 FROM cases
					 WHERE id=$1;`

	err := c.db.GetContext(ctx, &caseData, caseGetQuery, caseID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Case{}, customErrors.NoRowsCaseErr
		default:
			return models.Case{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
	}

	return caseData, nil
}

func (c caseRepo) GetCasesByLevel(ctx context.Context, level, cursor int) (models.CaseCursor, error) {
	var cases []models.Case
	var nextCursor null.Int
	var casesWithCursor models.CaseCursor

	// TODO: запрашивать нужно только те кейсы, которые он еще не оценил

	casesGetQueue := `SELECT id, camera_id, transport, violation_id, violation_value, level, datetime, photo_url
					  FROM cases WHERE level = $1 AND id > $2
					  ORDER BY id LIMIT $3;`

	err := c.db.SelectContext(ctx, &cases, casesGetQueue, level, cursor, c.CasesPerRequest+1)
	if err != nil {
		return models.CaseCursor{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}

	if len(cases) == c.CasesPerRequest+1 {
		nextCursor = null.IntFrom(int64(cases[len(cases)-1].ID))
		cases = cases[:len(cases)-1]
	}

	casesWithCursor.Cases = cases
	casesWithCursor.Cursor = nextCursor

	return casesWithCursor, nil
}

func (c caseRepo) DeleteCase(ctx context.Context, caseID int) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	caseDeleteQuery := `DELETE FROM cases WHERE id=$1;`

	res, err := tx.ExecContext(ctx, caseDeleteQuery, caseID)
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

func (c caseRepo) CreateRated(ctx context.Context, rated models.RatedBase) (int, error) {
	var createdRatedID int

	tx, err := c.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	caseCreateQuery := `INSERT INTO rated_cases (specialist_id, case_id, choice, date, status)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING id;`

	err = tx.QueryRowxContext(ctx, caseCreateQuery,
		rated.SpecialistID, rated.CaseID, rated.Choice, rated.Date, rated.Status).Scan(&createdRatedID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, customErrors.UniqueRatedErr
		}

		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdRatedID, nil
}

func (c caseRepo) GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error) {
	var rated []models.Rated
	var nextCursor null.Int
	var casesWithCursor models.RatedCursor

	casesGetQueue := `SELECT DISTINCT ON (case_id) id, specialist_id, case_id, choice, date, status
					  FROM rated_cases
					  WHERE status != 'Unknown' AND id > $1
					  ORDER BY case_id, id
					  LIMIT $2;`

	err := c.db.SelectContext(ctx, &rated, casesGetQueue, cursor, c.CasesPerRequest+1)
	if err != nil {
		return models.RatedCursor{}, err
	}

	if len(rated) == c.CasesPerRequest+1 {
		nextCursor = null.IntFrom(int64(rated[len(rated)-1].ID))
		rated = rated[:len(rated)-1]
	}

	casesWithCursor.Rated = rated
	casesWithCursor.Cursor = nextCursor

	return casesWithCursor, nil
}

func (c caseRepo) UpdateRatedStatus(ctx context.Context, newRated models.RatedUpdate) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	ratedUpdateStatusQuery := `UPDATE rated_cases SET status = $1 WHERE id = $2;`

	res, err := tx.ExecContext(ctx, ratedUpdateStatusQuery, newRated.Status, newRated.CaseID)
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
