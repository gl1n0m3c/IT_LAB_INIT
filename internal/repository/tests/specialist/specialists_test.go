package test_specialist

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository/tests"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
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

func TestCreateGetUpdateDeleteSpecialist(t *testing.T) {
	var createdIDs []int

	specRepo := repository.InitSpecialistsRepo(db)
	ctx := context.Background()

	// CreateCase tests
	for _, specialistCase := range testcaseSpecialistCreate {
		id, err := specRepo.Create(ctx, specialistCase)
		if err != nil {
			t.Errorf("CreateCase error: %v", err)
			continue
		}

		createdIDs = append(createdIDs, id)
	}

	// GetByID tests
	for i, id := range createdIDs {
		specialist, err := specRepo.GetByID(ctx, id)
		if err != nil {
			t.Errorf("GetByID error: %v", err)
			continue
		}

		assert.Equal(t, testcaseSpecialistCreate[i].Login, specialist.Login, "Compare Login")
		assert.Equal(t, testcaseSpecialistCreate[i].Fullname, specialist.Fullname, "Compare Fullname")
		assert.Equal(t, testcaseSpecialistCreate[i].PhotoUrl, specialist.PhotoUrl, "Compare PhotoUrl")
		assert.Equal(t, false, specialist.IsVerified, "Compare IsVerified")
		assert.Equal(t, 1, specialist.Level, "Compare Level")
		assert.Equal(t, true, utils.ComparePassword(specialist.Password, testcaseSpecialistCreate[i].Password), "Compare Password")
	}

	// GetByLogin tests
	for i, spec := range testcaseSpecialistCreate {
		specialist, err := specRepo.GetByLogin(ctx, spec.Login)
		if err != nil {
			t.Errorf("GetByID error: %v", err)
			continue
		}

		assert.Equal(t, testcaseSpecialistCreate[i].Login, specialist.Login, "Compare Login")
		assert.Equal(t, testcaseSpecialistCreate[i].Fullname, specialist.Fullname, "Compare Fullname")
		assert.Equal(t, testcaseSpecialistCreate[i].PhotoUrl, specialist.PhotoUrl, "Compare PhotoUrl")
		assert.Equal(t, false, specialist.IsVerified, "Compare IsVerified")
		assert.Equal(t, 1, specialist.Level, "Compare Level")
		assert.Equal(t, true, utils.ComparePassword(specialist.Password, testcaseSpecialistCreate[i].Password), "Compare Password")
	}

	// UpdateMain tests
	for i, id := range createdIDs {
		testcaseSpecialistUpdate[i].ID = id

		err := specRepo.UpdateMain(ctx, testcaseSpecialistUpdate[i], true)
		if err != nil {
			t.Errorf("UpdateMain error: %v", err)
			continue
		}

		specialist, err := specRepo.GetByID(ctx, id)
		if err != nil {
			t.Errorf("GetByID error: %v", err)
			continue
		}

		assert.Equal(t, testcaseSpecialistUpdate[i].ID, specialist.ID, "ID should match")
		assert.Equal(t, testcaseSpecialistUpdate[i].Login, specialist.Login, "Login should match")
		assert.Equal(t, true, utils.ComparePassword(specialist.Password, testcaseSpecialistUpdate[i].Password), "Password should match")
		assert.Equal(t, testcaseSpecialistUpdate[i].Fullname, specialist.Fullname, "Fullname should match")
		assert.Equal(t, testcaseSpecialistUpdate[i].PhotoUrl, specialist.PhotoUrl, "PhotoUrl should match")
		assert.Equal(t, testcaseSpecialistUpdate[i].Level, specialist.Level, "Level should match")
		assert.Equal(t, testcaseSpecialistUpdate[i].IsVerified, specialist.IsVerified, "IsVerified should match")

	}

	// DeleteCase tests
	for _, id := range createdIDs {
		err := specRepo.Delete(ctx, id)
		if err != nil {
			t.Errorf("DeleteCase error: %v", err)
			continue
		}
	}

	// Проверка на то, что записей не осталось
	for _, id := range createdIDs {
		_, err := specRepo.GetByID(ctx, id)
		assert.NotNil(t, err, "Expected an error when trying to GetByID a deleted specialist with ID %d, but got nil", id)
	}
}
