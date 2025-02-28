package models

import (
	"database/sql"
	"testing"

	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestConvertToBank(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected Bank
	}{
		{
			name: "Convert GetBankBySwiftCodeWithCountryRow",
			input: db.GetBankBySwiftCodeWithCountryRow{
				BankAddress: sql.NullString{String: "Test Address", Valid: true},
				BankName:    "Test Bank",
				CountryCode: "PL",
				CountryName: "Poland",
				BankType:    db.BankTypeHeadquarter,
				SwiftCode:   "TESTPLPWXXX",
			},
			expected: Bank{
				Address:       "Test Address",
				BankName:      "Test Bank",
				CountryCode:   "PL",
				CountryName:   "Poland",
				IsHeadquarter: true,
				SwiftCode:     "TESTPLPWXXX",
			},
		},
		{
			name: "Convert GetBanksBranchesBySwiftCodePrefixRow",
			input: db.GetBanksBranchesBySwiftCodePrefixRow{
				BankAddress: sql.NullString{String: "Branch Address", Valid: true},
				BankName:    "Branch Bank",
				CountryCode: "DE",
				BankType:    db.BankTypeBranch,
				SwiftCode:   "BRANCHDEXXX",
			},
			expected: Bank{
				Address:       "Branch Address",
				BankName:      "Branch Bank",
				CountryCode:   "DE",
				IsHeadquarter: false,
				SwiftCode:     "BRANCHDEXXX",
			},
		},
		{
			name: "Convert GetBanksByCountryCodeRow",
			input: db.GetBanksByCountryCodeRow{
				BankAddress: sql.NullString{String: "Country Bank Address", Valid: true},
				BankName:    "Country Bank",
				CountryCode: "US",
				BankType:    db.BankTypeBranch,
				SwiftCode:   "COUNTRYUSXXX",
			},
			expected: Bank{
				Address:       "Country Bank Address",
				BankName:      "Country Bank",
				CountryCode:   "US",
				IsHeadquarter: false,
				SwiftCode:     "COUNTRYUSXXX",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertToBank(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertToBanks(t *testing.T) {
	input := []db.GetBanksByCountryCodeRow{
		{
			BankAddress: sql.NullString{String: "Address 1", Valid: true},
			BankName:    "Bank 1",
			CountryCode: "PL",
			BankType:    db.BankTypeBranch,
			SwiftCode:   "SWIFT1",
		},
		{
			BankAddress: sql.NullString{String: "Address 2", Valid: true},
			BankName:    "Bank 2",
			CountryCode: "DE",
			BankType:    db.BankTypeHeadquarter,
			SwiftCode:   "SWIFT2",
		},
	}

	expected := []Bank{
		{
			Address:       "Address 1",
			BankName:      "Bank 1",
			CountryCode:   "PL",
			IsHeadquarter: false,
			SwiftCode:     "SWIFT1",
		},
		{
			Address:       "Address 2",
			BankName:      "Bank 2",
			CountryCode:   "DE",
			IsHeadquarter: true,
			SwiftCode:     "SWIFT2",
		},
	}

	result := ConvertToBanks(input)
	assert.Equal(t, expected, result)
}

func TestBankType(t *testing.T) {
	assert.Equal(t, db.BankTypeHeadquarter, BankType(true))
	assert.Equal(t, db.BankTypeBranch, BankType(false))
}
