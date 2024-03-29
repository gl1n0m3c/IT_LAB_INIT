package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"time"
)

type specialistsRepo struct {
	db                    *sqlx.DB
	j                     int
	specialistsPerRequest int
}

func InitSpecialistsRepo(db *sqlx.DB) Specialists {
	return specialistsRepo{db: db, j: viper.GetInt(config.J), specialistsPerRequest: viper.GetInt(config.EntitiesPerRequest)}
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

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, customErrors.UniqueSpecialistErr
		}

		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return createdSpecialistID, nil
}

func (s specialistsRepo) GetByID(ctx context.Context, specialistID int) (models.Specialist, error) {
	var specialist models.Specialist

	specialistGetQuery := `SELECT id, login, hashed_password, fullname, level, photo_url, is_verified, row
						FROM specialists
						WHERE id=$1;`

	err := s.db.GetContext(ctx, &specialist, specialistGetQuery, specialistID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Specialist{}, customErrors.NoRowsSpecialistIDErr
		default:
			return models.Specialist{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
	}

	return specialist, nil
}

func (s specialistsRepo) GetByLogin(ctx context.Context, specialistLogin string) (models.Specialist, error) {
	var specialist models.Specialist

	specialistGetQuery := `SELECT id, login, hashed_password, fullname, level, photo_url, is_verified, row
						FROM specialists
						WHERE login=$1;`

	err := s.db.GetContext(ctx, &specialist, specialistGetQuery, specialistLogin)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Specialist{}, customErrors.NoRowsSpecialistLoginErr
		default:
			return models.Specialist{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
	}

	return specialist, nil
}

func (s specialistsRepo) GetSpecialistRating(ctx context.Context, timeStart, timeEnd time.Time, cursor int) (models.RatingSpecialistCountCursor, error) {
	var (
		specialistsCursor models.RatingSpecialistCountCursor
		specialists       []models.RatingSpecialistCount
		nextCursor        null.Int
	)

	getRatingQuery := `SELECT s.id, s.fullname, s.level, s.photo_url, row,
					          COUNT(rc.id) AS total_cases,
					          COUNT(CASE WHEN rc.status = 'Correct' THEN 1 END) AS correct_cases,
					          COUNT(CASE WHEN rc.status = 'Unknown' THEN 1 END) AS unknown_cases
					   FROM specialists s
					   LEFT JOIN rated_cases rc ON s.id = rc.specialist_id AND rc.datetime BETWEEN $1 AND $2
					   WHERE s.id >= $3
					   GROUP BY s.id, s.login, s.fullname, s.level, s.photo_url, s.is_verified
					   LIMIT $4;`

	rows, err := s.db.QueryxContext(ctx, getRatingQuery, timeStart, timeEnd, cursor, s.specialistsPerRequest+1)
	if err != nil {
		return models.RatingSpecialistCountCursor{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}

	for rows.Next() {
		var ratingSpecialist models.RatingSpecialistCount

		err := rows.Scan(&ratingSpecialist.ID, &ratingSpecialist.Fullname, &ratingSpecialist.Level, &ratingSpecialist.PhotoUrl, &ratingSpecialist.Row,
			&ratingSpecialist.Total, &ratingSpecialist.Correct, &ratingSpecialist.Unknown)
		if err != nil {
			return models.RatingSpecialistCountCursor{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
		}

		specialists = append(specialists, ratingSpecialist)
	}

	if len(specialists) == s.specialistsPerRequest+1 {
		nextCursor = null.IntFrom(int64(specialists[len(specialists)-1].ID))
		specialists = specialists[:len(specialists)-1]
	}

	specialistsCursor.Specialists = specialists
	specialistsCursor.Cursor = nextCursor

	return specialistsCursor, nil
}

func (s specialistsRepo) GetFulRating(ctx context.Context) ([]models.RatingSpecialistFul, error) {
	var specialistsRating []models.RatingSpecialistFul

	getRatingQuery := `SELECT s.id, s.level, s.fullname,
					       COUNT(CASE WHEN rc.status = 'Correct' THEN 1 END) * 1.0 / NULLIF(COUNT(rc.id), 0) AS rating
					   FROM specialists s
					   LEFT JOIN rated_cases rc ON s.id = rc.specialist_id
					   GROUP BY s.id, s.level, s.fullname
					   ORDER BY rating DESC, COUNT(CASE WHEN rc.status = 'Correct' THEN 1 END);`

	rows, err := s.db.QueryxContext(ctx, getRatingQuery)
	if err != nil {
		return []models.RatingSpecialistFul{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var specialistRating models.RatingSpecialistFul

		err := rows.Scan(&specialistRating.SpecialistID, &specialistRating.Level, &specialistRating.Fullname, &specialistRating.Rating)
		if err != nil {
			return []models.RatingSpecialistFul{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}

		specialistsRating = append(specialistsRating, specialistRating)
	}

	if err := rows.Err(); err != nil {
		return []models.RatingSpecialistFul{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	return specialistsRating, nil
}

func (s specialistsRepo) GetOnlyRating(timeStart, timeEnd time.Time) ([]models.RatingSpecialistID, error) {
	var ratings []models.RatingSpecialistID

	getRatingQuery := `SELECT s.id, s.level, COUNT(CASE WHEN rc.status = 'Correct' THEN 1 END) * 1.0 / COUNT(rc.id) AS rating
					   FROM specialists s
					   LEFT JOIN rated_cases rc ON s.id = rc.specialist_id AND rc.datetime BETWEEN $1 AND $2
					   GROUP BY s.id, s.level
					   HAVING COUNT(rc.id) >= $3
					   ORDER BY COUNT(CASE WHEN rc.status = 'Correct' THEN 1 END) * 1.0 / COUNT(rc.id) DESC, COUNT(CASE WHEN rc.status = 'Correct' THEN 1 END);`

	rows, err := s.db.Queryx(getRatingQuery, timeStart, timeEnd, s.j)
	if err != nil {
		return []models.RatingSpecialistID{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.QueryRrr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var rating models.RatingSpecialistID

		err := rows.Scan(&rating.ID, &rating.Level, &rating.Rating)
		if err != nil {
			return []models.RatingSpecialistID{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}

		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		return []models.RatingSpecialistID{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	return ratings, nil
}

func (s specialistsRepo) UpdateSpecialistsIncDecLevel(incrementIDs, decrementIDs []int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	query1 := `UPDATE specialists SET level = level + 1 WHERE id = ANY($1)`
	query2 := `UPDATE specialists SET level = level - 1 WHERE id = ANY($1)`

	res, err := s.db.Exec(query1, pq.Array(incrementIDs))
	if err != nil {
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

	if count != int64(len(incrementIDs)) {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	res, err = s.db.Exec(query2, pq.Array(decrementIDs))
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}
	count, err = res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	if count != int64(len(decrementIDs)) {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	return nil
}

func (s specialistsRepo) UpdateMain(ctx context.Context, specialistUpdate models.Specialist, newPasswordFlag bool) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	if newPasswordFlag {
		specialistUpdate.Password = string(utils.HashPassword(specialistUpdate.Password))
	}

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
