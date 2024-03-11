package utils

import (
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
)

func LoadFixtures(db *sqlx.DB) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Paths(
			"../internal/fixtures/cameras.yml",
			"../internal/fixtures/specialists.yml",
			"../internal/fixtures/violations.yml",
			"../internal/fixtures/contacts.yml",
			"../internal/fixtures/cases.yml",
			"../internal/fixtures/rated_cases.yml",
		),
	)

	if err != nil {
		panic(err)
	}

	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

func ClearDatabase(db *sqlx.DB) {
	tables := []string{"cases", "contacts", "violations", "specialists", "cameras"}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}
}
