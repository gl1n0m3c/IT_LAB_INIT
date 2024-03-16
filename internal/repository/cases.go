package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

type caseRepo struct {
	db              *sqlx.DB
	casesPerRequest int
}

func InitCaseRepo(
	db *sqlx.DB,
) Cases {
	return caseRepo{
		db:              db,
		casesPerRequest: viper.GetInt(config.EntitiesPerRequest),
	}
}

func (c caseRepo) CreateCase(ctx context.Context, caseData models.CaseBase) (int, error) {
	var createdCaseID int

	tx, err := c.db.Beginx()
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	caseCreateQuery := `INSERT INTO cases (camera_id, transport, violation_id, violation_value, level, current_level, datetime, photo_url)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
						RETURNING id;`

	err = tx.QueryRowxContext(ctx, caseCreateQuery,
		caseData.CameraID, caseData.Transport, caseData.ViolationID, caseData.ViolationValue,
		caseData.Level, caseData.Level, caseData.Datetime, caseData.PhotoUrl).Scan(&createdCaseID)
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

func (c caseRepo) UpdateCaseLevel(ctx context.Context, caseID, level int) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE cases SET current_level = $1 WHERE id=$2;`

	res, err := c.db.ExecContext(ctx, updateQuery, level, caseID)
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

func (c caseRepo) UpdateCaseSetSolved(ctx context.Context, caseID int, rightChoice bool) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	updateCaseSolvedQuery := `UPDATE cases SET is_solved=true WHERE id=$1;`

	updateRatedSpecialistsQuery := `UPDATE rated_cases
									SET status =
										CASE
									WHEN choice = $1 THEN 'Correct'::status_type
									ELSE 'Incorrect'::status_type
									END
									WHERE case_id = $2;`

	updateSpecialistsQuery := `UPDATE specialists s
        					   SET current_row = CASE
							   	WHEN rc.status = 'Correct'::status_type THEN s.current_row + 1
							   	ELSE 0
							   END,
							   row = CASE
							   	WHEN rc.status = 'Correct'::status_type AND s.current_row + 1 > s.row THEN s.current_row + 1
							   	ELSE s.row
							   END
        					   FROM rated_cases rc
        					   WHERE rc.case_id = $1 AND rc.specialist_id = s.id
        					   RETURNING s.id`

	res, err := c.db.ExecContext(ctx, updateCaseSolvedQuery, caseID)
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

	_, err = c.db.ExecContext(ctx, updateRatedSpecialistsQuery, rightChoice, caseID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}

	_, err = c.db.ExecContext(ctx, updateSpecialistsQuery, caseID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}

	return nil
}

func (c caseRepo) GetFineData(ctx context.Context, caseID int) (models.FineData, error) {
	var fineData models.FineData

	getFineDataQuery := `SELECT cn.contacts, c.photo_url, cm.coordinates, c.violation_value, v.type, v.amount, c.datetime
						 FROM cases c
						 JOIN violations v ON c.violation_id = v.id
						 JOIN contacts cn ON c.transport = cn.transport
						 JOIN cameras cm ON c.camera_id = cm.id
						 WHERE c.id=$1`
	err := c.db.QueryRowxContext(ctx, getFineDataQuery, caseID).Scan(&fineData.Mail, &fineData.PhotoUrl, &fineData.Coordinated,
		&fineData.ViolationValue, &fineData.Violation.Type, &fineData.Violation.Amount, &fineData.Date)
	if err != nil {
		return models.FineData{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	var contacts map[string]string
	if err := json.Unmarshal([]byte(fineData.Mail), &contacts); err != nil {
		return models.FineData{}, err // Используйте вашу собственную обработку ошибок
	}
	fineData.Mail = contacts["email"]

	return fineData, nil
}

func (c caseRepo) GetCaseLevelSolvedRatingsTrueByID(ctx context.Context, caseID, specialistLevel int) (int, int, int, bool, error) {
	var level, num, numTrue int
	var isSolved bool

	caseGetQuery := `SELECT c.current_level, c.is_solved, COUNT(rs.case_id), SUM(CASE WHEN rs.choice = true THEN 1 ELSE 0 END)
					 FROM cases c
					 LEFT JOIN rated_cases rs ON c.id = rs.case_id
					 LEFT JOIN specialists s ON rs.specialist_id = s.id
					 WHERE c.id = $1 AND s.level = $2
					 GROUP BY c.current_level, c.is_solved;`

	err := c.db.QueryRowxContext(ctx, caseGetQuery, caseID, specialistLevel).Scan(&level, &isSolved, &num, &numTrue)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, 0, 0, false, customErrors.UserBadLevel
		default:
			return 0, 0, 0, false, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
	}

	return level, num, numTrue, isSolved, nil
}

func (c caseRepo) GetCasesByLevel(ctx context.Context, specialistID, level, cursor int) (models.CaseCursor, error) {
	var cases []models.CaseViolations
	var nextCursor null.Int
	var casesWithCursor models.CaseCursor

	casesGetQueue := `SELECT c.id, c.camera_id, c.transport, c.violation_id, c.violation_value, c.level, c.datetime, c.photo_url,
					  v.type, v.amount
					  FROM cases c
					  LEFT JOIN violations v ON c.violation_id = v.id
					  LEFT JOIN rated_cases rc ON c.id = rc.case_id AND rc.specialist_id = $1
					  WHERE c.level = $2 AND rc.id IS NULL AND c.id >= $3 AND c.is_solved = false
					  ORDER BY id LIMIT $4;`

	err := c.db.SelectContext(ctx, &cases, casesGetQueue, specialistID, level, cursor, c.casesPerRequest+1)
	if err != nil {
		return models.CaseCursor{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}

	if len(cases) == c.casesPerRequest+1 {
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

	caseCreateQuery := `INSERT INTO rated_cases (specialist_id, case_id, choice, datetime, status)
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

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, customErrors.UniqueRatedErr
		}

		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdRatedID, nil
}

func (c caseRepo) GetNumberRatedByCaseID(ctx context.Context, caseID int) (int, error) {
	var number int

	getNumberQuery := `SELECT COUNT(*) FROM rated_cases WHERE case_id=$1;`

	err := c.db.QueryRowxContext(ctx, getNumberQuery, caseID).Scan(&number)
	if err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	return number, nil
}

func (c caseRepo) GetRatedSolved(ctx context.Context, cursor int) (models.RatedCursor, error) {
	var rated []models.Rated
	var nextCursor null.Int
	var casesWithCursor models.RatedCursor

	casesGetQueue := `SELECT DISTINCT ON (rc.case_id) rc.id, rc.specialist_id, rc.case_id, rc.choice, rc.datetime, rc.status,
					  c.camera_id, c.violation_value, v.type, v.amount, c.level, c.photo_url
					  FROM rated_cases rc
					  LEFT JOIN cases c ON rc.case_id = c.id
					  LEFT JOIN violations v ON c.violation_id = v.id
					  WHERE status != 'Unknown' AND rc.id >= $1 AND c.is_solved = true
					  ORDER BY case_id, id
					  LIMIT $2;`

	rows, err := c.db.QueryxContext(ctx, casesGetQueue, cursor, c.casesPerRequest+1)
	if err != nil {
		return models.RatedCursor{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var rt models.Rated

		err := rows.Scan(&rt.ID, &rt.SpecialistID, &rt.CaseID, &rt.Choice, &rt.Date, &rt.Status, &rt.CameraID,
			&rt.ViolationValue, &rt.Violation.Type, &rt.Violation.Amount, &rt.Level, &rt.PhotoUrl)
		if err != nil {
			return models.RatedCursor{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}

		rated = append(rated, rt)
	}

	if len(rated) == c.casesPerRequest+1 {
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

func (c caseRepo) GetFulCaseByID(ctx context.Context, caseID int) (models.CaseFul, error) {
	var caseFul models.CaseFul
	var ratedNum int

	caseGetQuery := `SELECT c.camera_id, c.transport, v.id, v.type, v.amount, c.violation_value, c.level, c.current_level, c.datetime, c.photo_url, c.is_solved, COUNT(*)
					 FROM cases c
					 LEFT JOIN violations v ON c.violation_id = v.id
					 LEFT JOIN rated_cases rc ON c.id = rc.case_id
					 WHERE c.id = $1
					 GROUP BY c.camera_id, c.transport, v.id, v.type, v.amount, c.violation_value, c.level, c.current_level, c.datetime, c.photo_url, c.is_solved`

	err := c.db.QueryRowxContext(ctx, caseGetQuery, caseID).Scan(&caseFul.CameraID, &caseFul.Transport, &caseFul.ViolationID, &caseFul.Violation.Type,
		&caseFul.Violation.Amount, &caseFul.ViolationValue,
		&caseFul.Level, &caseFul.CurrentLevel, &caseFul.Datetime, &caseFul.PhotoUrl, &caseFul.IsSolved, &ratedNum)
	if err != nil {
		return models.CaseFul{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if ratedNum > 1 {
		ratedGetQuery := `SELECT rc.id, rc.status, rc.datetime, s.id, s.fullname, s.level, s.photo_url
					  	  FROM rated_cases rc
					  	  LEFT JOIN specialists s ON rc.specialist_id = s.id
					  	  WHERE rc.case_id = $1`

		rows, err := c.db.QueryxContext(ctx, ratedGetQuery, caseID)
		defer rows.Close()

		var ratedCovers []models.RatedCover

		for rows.Next() {
			var ratedCover models.RatedCover

			err := rows.Scan(&ratedCover.ID, &ratedCover.Status, &ratedCover.Date,
				&ratedCover.SpecialistCover.ID, &ratedCover.SpecialistCover.Fullname, &ratedCover.SpecialistCover.Level, &ratedCover.SpecialistCover.PhotoUrl)
			if err != nil {
				return models.CaseFul{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
			}

			ratedCovers = append(ratedCovers, ratedCover)
		}

		err = rows.Err()
		if err != nil {
			return models.CaseFul{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
		}

		caseFul.RatedCovers = &ratedCovers
	} else {
		caseFul.RatedCovers = nil
	}

	return caseFul, nil

}
