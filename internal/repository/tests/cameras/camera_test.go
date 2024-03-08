package cameras

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository/tests"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error

	tests.InitTestConfig()

	connectionString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		viper.GetString(tests.TestDBUser),
		viper.GetString(tests.TestDBPassword),
		viper.GetString(tests.TestDBHost),
		viper.GetInt(tests.TestDBPort),
		viper.GetString(tests.TestDBName),
	)

	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Could not connect to the tests database: %v", err)
	}

	code := m.Run()

	err = db.Close()
	if err != nil {
		return
	}

	os.Exit(code)
}

func TestCreateGetDeleteCameras(t *testing.T) {
	var createdIDs []string

	cameraRepo := repository.InitCameraRepo(db)
	ctx := context.Background()

	// CreateCase tests
	for _, cameraCase := range testcaseCameraCreate {
		id, err := cameraRepo.Create(ctx, cameraCase)
		if err != nil {
			t.Errorf("CreateCase error: %v", err)
			continue
		}

		createdIDs = append(createdIDs, id)
	}

	// Get tests
	for i, id := range createdIDs {
		camera, err := cameraRepo.Get(ctx, id)
		if err != nil {
			t.Errorf("Get error: %v", err)
			continue
		}

		assert.Equal(t, testcaseCameraCreate[i].Type, camera.Type, "Compare Type")
		assert.Equal(t, testcaseCameraCreate[i].Coordinates, camera.Coordinates, "Compare Coordinates")
		assert.Equal(t, testcaseCameraCreate[i].Description, camera.Description, "Compare PhotoUrl")
	}

	// DeleteCase tests
	for _, id := range createdIDs {
		err := cameraRepo.Delete(ctx, id)
		if err != nil {
			t.Errorf("DeleteCase error: %v", err)
			continue
		}
	}

	// Проверка на то, что записей не осталось
	for _, id := range createdIDs {
		_, err := cameraRepo.Get(ctx, id)
		assert.NotNil(t, err, "Expected an error when trying to Get a deleted specialist with ID %d, but got nil", id)
	}
}
