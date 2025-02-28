package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	models "github.com/mateuszkochelski/SwiftCodeDb/models"
	store "github.com/mateuszkochelski/SwiftCodeDb/repository"
)

type BankHandler struct {
	queries *db.Queries
}

func NewBankHandler(dataBase *sql.DB) *BankHandler {
	return &BankHandler{queries: db.New(dataBase)}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type BankResponse struct {
	Address       string        `json:"address"`
	BankName      string        `json:"bankName"`
	CountryISO2   string        `json:"countryISO2"`
	CountryName   string        `json:"countryName"`
	IsHeadquarter bool          `json:"isHeadquarter"`
	SwiftCode     string        `json:"swiftCode"`
	Branches      []models.Bank `json:"branches,omitempty"`
}

type CountryBanksResponse struct {
	CountryCode string        `json:"countryISO2"`
	CountryName string        `json:"countryName"`
	SwiftCodes  []models.Bank `json:"swiftCodes"`
}

func sendJSONError(w http.ResponseWriter, statusCode int, errMsg string, details ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error: fmt.Sprintf(errMsg, details...),
	}

	json.NewEncoder(w).Encode(response)
}

func (h *BankHandler) GetBanksBySwiftCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var response any
	swiftCode := strings.TrimPrefix(r.URL.Path, "/v1/swift-codes/")
	bankQueryResult, err := h.queries.GetBankBySwiftCodeWithCountry(r.Context(), swiftCode)
	if err != nil {
		sendJSONError(w, http.StatusNotFound, "Not Found: ")
		return
	}

	bank := models.ConvertToBank(bankQueryResult)

	if bank.IsHeadquarter {
		swiftCodePrefix, ok := strings.CutSuffix(swiftCode, "XXX")
		if !ok {
			sendJSONError(w, http.StatusInternalServerError, "Data inconsistency: headquarter bank must have SWIFT code ending in 'XXX'")
			return
		}
		sqlRegex := swiftCodePrefix + "___"
		queryParams := db.GetBanksBranchesBySwiftCodePrefixParams{
			SwiftCode:   sqlRegex,
			SwiftCode_2: swiftCode,
		}

		banksQueryResult, err := h.queries.GetBanksBranchesBySwiftCodePrefix(r.Context(), queryParams)
		if err != nil && err != sql.ErrNoRows {
			sendJSONError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		branches := models.ConvertToBanks(banksQueryResult)

		response = BankResponse{
			Address:       bank.Address,
			BankName:      bank.BankName,
			CountryISO2:   bank.CountryCode,
			CountryName:   bank.CountryName,
			IsHeadquarter: bank.IsHeadquarter,
			SwiftCode:     bank.SwiftCode,
			Branches:      branches,
		}
	} else {
		response = bank
		_, ok := strings.CutSuffix(swiftCode, "XXX")
		if ok {
			http.Error(w, "Data inconsistency: branch bank should not have SWIFT code ending in 'XXX'", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BankHandler) GetBanksByContryCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	countryCode := strings.TrimPrefix(r.URL.Path, "/v1/swift-codes/country/")
	country, err := h.queries.GetCountry(r.Context(), countryCode)
	if err != nil {
		sendJSONError(w, http.StatusNotFound, "Not found")
	}

	banks, err := h.queries.GetBanksByCountryCode(r.Context(), countryCode)
	if err != nil {
		sendJSONError(w, http.StatusNotFound, "Not found")
	}

	response := CountryBanksResponse{
		CountryCode: country.CountryCode,
		CountryName: country.CountryName,
	}
	for _, bank := range banks {
		response.SwiftCodes = append(response.SwiftCodes, models.ConvertToBank(bank))
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BankHandler) CreateBank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var request models.Bank
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	_, ok := strings.CutSuffix(request.SwiftCode, "XXX")
	if request.IsHeadquarter && !ok {
		sendJSONError(w, http.StatusUnprocessableEntity, "Headquarter swift code should end with XXX")
		return
	}
	if !request.IsHeadquarter && ok {
		sendJSONError(w, http.StatusUnprocessableEntity, "Branch swift code shouldnt end with XXX")
		return
	}

	country := db.CreateCountryParams{
		CountryCode: request.CountryCode,
		CountryName: request.CountryName,
	}
	err = store.InsertCountryWithValidation(h.queries, country)
	if err != nil {
		sendJSONError(w, http.StatusUnprocessableEntity, "Error during country insertion: %s", err.Error())
		return
	}

	bank := db.CreateBankParams{
		BankName:    request.BankName,
		SwiftCode:   request.SwiftCode,
		BankAddress: sql.NullString{String: request.Address, Valid: len(request.Address) != 0},
		CountryCode: request.CountryCode,
		BankType:    models.BankType(request.IsHeadquarter),
	}
	err = store.InsertBankWithValidation(h.queries, bank)
	if err != nil {
		sendJSONError(w, http.StatusUnprocessableEntity, "Error during bank insertion: %s,%s", err.Error(), bank.CountryCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bank created successfully"})
}

func (h *BankHandler) DeleteBank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendJSONError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	swiftCode := strings.TrimPrefix(r.URL.Path, "/v1/swift-codes/")

	err := h.queries.DeleteBankBySwiftCode(r.Context(), swiftCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendJSONError(w, http.StatusNotFound, "Bank not found")
			return
		}

		sendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bank deleted successfully"})
}

func (h *BankHandler) HandleSwiftCodes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetBanksBySwiftCode(w, r)
	case http.MethodDelete:
		h.DeleteBank(w, r)
	default:
		sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
