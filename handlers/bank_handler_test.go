package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	models "github.com/mateuszkochelski/SwiftCodeDb/models"
	store "github.com/mateuszkochelski/SwiftCodeDb/repository"
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
	_, err = queries.DeleteBankBySwiftCode(context.Background(), "ABCABCABXXX")
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

func TestCreateAndDeleteBank(t *testing.T) {
	dbConn := setupTestDB()
	defer dbConn.Close()

	handler := NewBankHandler(dbConn)

	// Step 1: Create a bank
	bank := models.Bank{
		BankName:      "Test Bank",
		SwiftCode:     "TESTBANKXXX",
		Address:       "Test Address",
		CountryCode:   "CC",
		CountryName:   "CCCC",
		IsHeadquarter: true,
	}

	bankJSON, err := json.Marshal(bank)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes/", strings.NewReader(string(bankJSON)))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	handler.CreateBank(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/TESTBANKXXX", nil)
	deleteResp := httptest.NewRecorder()
	handler.DeleteBank(deleteResp, deleteReq)
	assert.Equal(t, http.StatusCreated, deleteResp.Code)

	deleteReq2 := httptest.NewRequest(http.MethodDelete, "/v1/swift-codes/TESTBANKXXX", nil)
	deleteResp2 := httptest.NewRecorder()
	handler.DeleteBank(deleteResp2, deleteReq2)
	assert.Equal(t, http.StatusNotFound, deleteResp2.Code)
}

func TestGetBanksByCountryCode(t *testing.T) {
	dbConn := setupTestDB()
	defer dbConn.Close()

	handler := NewBankHandler(dbConn)

	country := db.CreateCountryParams{
		CountryCode: "PL",
		CountryName: "POLAND",
	}
	err := store.InsertCountryWithValidation(db.New(dbConn), country)
	assert.NoError(t, err)

	bank := db.CreateBankParams{
		BankName:    "Test Bank",
		SwiftCode:   "TESTBANKXXX",
		BankAddress: sql.NullString{String: "Test Address", Valid: true},
		CountryCode: "PL",
		BankType:    models.BankType(true),
	}
	err = store.InsertBankWithValidation(db.New(dbConn), bank)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/country/PL", nil)
	resp := httptest.NewRecorder()
	handler.GetBanksByContryCode(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response CountryBanksResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "PL", response.CountryCode)
	assert.Equal(t, "POLAND", response.CountryName)
	assert.NotEmpty(t, response.SwiftCodes)
}
