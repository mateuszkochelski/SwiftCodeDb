package models

import (
	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
)

type Bank struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryCode   string `json:"countryISO2"`
	CountryName   string `json:"countryName,omitempty"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func ConvertToBank[T any](row T) Bank {
	switch r := any(row).(type) {
	case db.GetBankBySwiftCodeWithCountryRow:
		return Bank{
			Address:       r.BankAddress.String,
			BankName:      r.BankName,
			CountryCode:   r.CountryCode,
			CountryName:   r.CountryName,
			IsHeadquarter: r.BankType == db.BankTypeHeadquarter,
			SwiftCode:     r.SwiftCode,
		}
	case db.GetBanksBranchesBySwiftCodePrefixRow:
		return Bank{
			Address:       r.BankAddress.String,
			BankName:      r.BankName,
			CountryCode:   r.CountryCode,
			IsHeadquarter: r.BankType == db.BankTypeHeadquarter,
			SwiftCode:     r.SwiftCode,
		}
	case db.GetBanksByCountryCodeRow:
		return Bank{
			Address:       r.BankAddress.String,
			BankName:      r.BankName,
			CountryCode:   r.CountryCode,
			IsHeadquarter: r.BankType == db.BankTypeHeadquarter,
			SwiftCode:     r.SwiftCode,
		}
	default:
		panic("Unsupported type for ConvertBank")
	}
}

func ConvertToBanks[T any](rows []T) []Bank {
	var banks []Bank
	for _, row := range rows {
		banks = append(banks, ConvertToBank(row))
	}

	return banks
}

func BankType(isHeadquarer bool) db.BankType {
	if isHeadquarer {
		return db.BankTypeHeadquarter
	}
	return db.BankTypeBranch
}
