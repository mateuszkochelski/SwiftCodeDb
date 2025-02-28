package repository

import (
	"database/sql"
	"fmt"
	"testing"

	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	"github.com/stretchr/testify/assert"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://test:test@postgresTestDB:5432/testdb?sslmode=disable"
)

func setupTestDB() (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("cannot connect to test database: %w", err)
	}

	return conn, nil
}

func TestInsertCountryWithValidation(t *testing.T) {
	testDB, err := setupTestDB()
	if err != nil {
		return
	}
	defer testDB.Close()
	queries := db.New(testDB)
	newCountry := db.CreateCountryParams{
		CountryCode: "PL",
		CountryName: "POLAND",
	}

	err = InsertCountryWithValidation(queries, newCountry)
	assert.NoError(t, err)

	err = InsertCountryWithValidation(queries, db.CreateCountryParams{
		CountryCode: "PL",
		CountryName: "GERMANY",
	})
	assert.ErrorContains(t, err, "inconsistency in request")
}
