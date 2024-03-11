package cases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository/tests"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
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

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Paths(
			"../../../fixtures/cameras.yml",
			"../../../fixtures/specialists.yml",
			"../../../fixtures/violations.yml",
			"../../../fixtures/contacts.yml",
		),
	)
	if err != nil {
		log.Fatalf("Error creating fixtures: %v", err)
	}

	if err := fixtures.Load(); err != nil {
		log.Fatalf("Error loading fixtures: %v", err)
	}

	code := m.Run()

	db.Exec("DELETE FROM cameras")
	db.Exec("DELETE FROM specialists")
	db.Exec("DELETE FROM violations")
	db.Exec("DELETE FROM contacts")

	err = db.Close()
	if err != nil {
		log.Fatalf("Could not close the tests database: %v", err)
	}

	os.Exit(code)
}

func TestCreateGetDeleteCasesRated(t *testing.T) {
	var createdIDs []int

	caseRepo := repository.InitCaseRepo(db, 2)

	// Creating
	for _, caseData := range testCases {
		ctx, cansel := context.WithTimeout(context.Background(), time.Second*2)
		defer cansel()

		id, err := caseRepo.CreateCase(ctx, caseData)
		if err != nil {
			t.Errorf(err.Error())
		}
		createdIDs = append(createdIDs, id)
	}

	ctx := context.Background()

	caseCursor, err := caseRepo.GetCasesByLevel(ctx, 2, 0)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: 0, Valid: false}}, caseCursor.Cursor)

	// Get
	for i, id := range createdIDs {
		caseCursor, err := caseRepo.GetCasesByLevel(ctx, 1, id-1)
		if err != nil {
			t.Errorf(err.Error())
		}

		if i > 1 {
			assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: 0, Valid: false}}, caseCursor.Cursor)
		} else {
			assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: int64(id + 2), Valid: true}}, caseCursor.Cursor)
		}
	}

	// Create rated
	var createdRatedIDs []int
	for i, rated := range testCasesRated {
		rated.CaseID = createdIDs[i]

		id, err := caseRepo.CreateRated(ctx, rated)
		if err != nil {
			t.Errorf(err.Error())
		}

		createdRatedIDs = append(createdRatedIDs, id)
	}

	// Get rated
	for i, id := range createdRatedIDs {
		if i == 3 {
			break
		}

		caseCursor, err := caseRepo.GetRatedSolved(ctx, id-1)
		if err != nil {
			t.Errorf(err.Error())
		}

		if i > 0 {
			assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: 0, Valid: false}}, caseCursor.Cursor)
		} else {
			assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: int64(id + 2), Valid: true}}, caseCursor.Cursor)
		}
	}

	// UpdateMain rated status
	for _, id := range createdRatedIDs {
		err := caseRepo.UpdateRatedStatus(ctx, models.RatedUpdate{
			CaseID: id,
			Status: "Correct",
		})
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	// Get rated
	for i, id := range createdRatedIDs {
		caseCursor, err := caseRepo.GetRatedSolved(ctx, id-1)
		if err != nil {
			t.Errorf(err.Error())
		}

		if i > 2 {
			assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: 0, Valid: false}}, caseCursor.Cursor)
		} else {
			assert.Equal(t, null.Int{NullInt64: sql.NullInt64{Int64: int64(id + 2), Valid: true}}, caseCursor.Cursor)
		}
	}

	// Delete
	for _, id := range createdIDs {
		err := caseRepo.DeleteCase(ctx, id)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

}
