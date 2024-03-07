package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/custom_errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type cameraRepo struct {
	db *sqlx.DB
}

func InitCameraRepo(db *sqlx.DB) Cameras {
	return cameraRepo{db: db}
}

func (c cameraRepo) Create(ctx context.Context, camera models.CameraBase) (string, error) {
	uuidBytes, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	key := uuidBytes.String()

	tx, err := c.db.Beginx()
	if err != nil {
		return "", utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	cameraQueue := `INSERT INTO cameras (id, type, description, coordinates)
					VALUES ($1, $2, $3, $4);`

	coordinates := fmt.Sprintf("%g,%g", camera.Coordinates[0], camera.Coordinates[1])

	res, err := tx.ExecContext(ctx, cameraQueue, key, camera.Type, camera.Description, coordinates)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return "", utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return "", utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}

	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return "", utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return "", utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return "", utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return key, nil
}

func (c cameraRepo) Get(ctx context.Context, cameraID string) (models.Camera, error) {
	var camera models.Camera
	var coords string
	var latitude, longitude float64

	cameraGetQuery := `SELECT id, type, description, coordinates
						FROM cameras
						WHERE id=$1;`

	err := c.db.QueryRowxContext(ctx, cameraGetQuery, cameraID).Scan(&camera.ID, &camera.Type, &camera.Description, &coords)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Camera{}, customErrors.NoRowsCameraErr
		default:
			return models.Camera{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
		}
	}

	_, _ = fmt.Sscanf(coords, "%g,%g", &latitude, &longitude)
	camera.Coordinates = [2]float64{latitude, longitude}

	return camera, nil
}

func (c cameraRepo) Delete(ctx context.Context, cameraID string) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	cameraDeleteQuery := `DELETE FROM cameras WHERE id=$1;`

	res, err := tx.ExecContext(ctx, cameraDeleteQuery, cameraID)
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
		return customErrors.NoRowsCameraErr
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}
