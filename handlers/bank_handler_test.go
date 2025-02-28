package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	"github.com/stretchr/testify/assert"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://test:test@postgresTestDB:5432/testdb?sslmode=disable"
)

func setupTestDB() *sql.DB {
	conn, _ := sql.Open(dbDriver, dbSource)
	err := conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	return conn
}

func Test_create_get_delete_succeed(t *testing.T) {
	testDB := setupTestDB()
	defer testDB.Close()

	bankHandler := NewBankHandler(testDB)

	requestBody := map[string]interface{}{
		"address":       "nowhere in poland",
		"bankName":      "Pekao",
		"countryISO2":   "PL",
		"countryName":   "POLAND",
		"isHeadquarter": true,
		"swiftCode":     "ABCABCABXXX",
	}

	jsonBody, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	postReq, err := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	postReq.Header.Set("Content-Type", "application/json")

	postRR := httptest.NewRecorder()
	postHandler := http.HandlerFunc(bankHandler.CreateBank)
	postHandler.ServeHTTP(postRR, postReq)

	assert.Equal(t, http.StatusCreated, postRR.Code, postRR.Body.String(), "Expected status 201 Created")

	getReq, err := http.NewRequest("GET", "/v1/swift-codes/ABCABCABXXX", nil)
	assert.NoError(t, err)

	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(bankHandler.GetBanksBySwiftCode)
	getHandler.ServeHTTP(getRR, getReq)

	assert.Equal(t, http.StatusOK, getRR.Code, "Expected status 200 OK")
	assert.Contains(t, getRR.Body.String(), `"Pekao"`, "Response should contain the bank name 'Pekao'")

	queries := db.New(testDB)
	err = queries.DeleteBankBySwiftCode(context.Background(), "ABCABCABXXX")
	assert.NoError(t, err)
}

func Test_create_fails_bad_country_code(t *testing.T) {
	testDB := setupTestDB()
	defer testDB.Close()

	bankHandler := NewBankHandler(testDB)

	requestBody := map[string]interface{}{
		"address":       "nowhere in poland",
		"bankName":      "Pekao",
		"countryISO2":   "Pl",
		"countryName":   "POLAND",
		"isHeadquarter": true,
		"swiftCode":     "ABCABCABXXX",
	}

	jsonBody, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	postReq, err := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	postReq.Header.Set("Content-Type", "application/json")

	postRR := httptest.NewRecorder()
	postHandler := http.HandlerFunc(bankHandler.CreateBank)
	postHandler.ServeHTTP(postRR, postReq)

	assert.Equal(t, http.StatusUnprocessableEntity, postRR.Code, postRR.Body.String(), "Expected uprocessable entity status 422 error code")
}
