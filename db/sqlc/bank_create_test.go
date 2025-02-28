package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_create_bank_succeded_when_swift_has_11_letters_country_uppercase_iso_code_uppercase(t *testing.T) {

	countryParams := CreateCountryParams{
		CountryCode: "PL",
		CountryName: "POLAND",
	}
	bankParams := CreateBankParams{
		SwiftCode:   "12345678XXX",
		BankName:    "Pekao",
		CountryCode: "PL",
		BankType:    BankTypeHeadquarter,
	}
	country, err := testQueries.CreateCountry(context.Background(), countryParams)
	require.NoError(t, err)
	require.NotEmpty(t, country)

	bank, err := testQueries.CreateBank(context.Background(), bankParams)
	require.NoError(t, err)
	require.NotEmpty(t, bank)
}

func Test_create_bank_error_when_swift_code_hasnt_11_letters(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678XXXX",
		BankName:    "Pekao",
		CountryCode: "PL",
		BankType:    BankTypeHeadquarter,
	}

	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, bank)
}
func Test_create_bank_error_when_country_code_not_uppercase(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678XXXX",
		BankName:    "Pekao",
		CountryCode: "Pl",
		BankType:    BankTypeHeadquarter,
	}

	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, bank)
}

func Test_create_bank_error_when_country_name_not_uppercase(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678XXXX",
		BankName:    "Pekao",
		CountryCode: "EN",
		BankType:    BankTypeHeadquarter,
	}

	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, bank)
}

func Test_create_bank_error_swift_code_not_ends_with_xxx_and_bank_type_headquarter(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678123",
		BankName:    "Pekao",
		CountryCode: "EN",
		BankType:    BankTypeHeadquarter,
	}

	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, bank)
}

func Test_create_bank_succeed_swift_code_ends_with_xxx_and_bank_type_headquarter(t *testing.T) {
	countryArg := CreateCountryParams{
		CountryCode: "EN",
		CountryName: "ENGLAND",
	}
	bankArg := CreateBankParams{
		SwiftCode:   "12345678XXX",
		BankName:    "Pekao",
		CountryCode: "EN",
		BankType:    BankTypeHeadquarter,
	}
	country, err := testQueries.CreateCountry(context.Background(), countryArg)
	require.NoError(t, err)
	require.NotEmpty(t, country)
	bank, err := testQueries.CreateBank(context.Background(), bankArg)
	require.NoError(t, err)
	require.NotEmpty(t, bank)
}

func Test_create_bank_succeed_swift_code_not_ends_with_xxx_and_bank_type_branch(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678ASD",
		BankName:    "Pekao",
		CountryCode: "EN",
		BankType:    BankTypeBranch,
	}

	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bank)
}

func Test_create_bank_error_swift_code_ends_with_xxx_and_bank_type_branch(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678XXX",
		BankName:    "Pekao",
		CountryCode: "EN",
		BankType:    BankTypeBranch,
	}

	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, bank)
}

func Test_create_bank_error_empty_bank_name(t *testing.T) {
	arg := CreateBankParams{
		SwiftCode:   "12345678XXX",
		CountryCode: "EN",
		BankType:    BankTypeHeadquarter,
	}
	bank, err := testQueries.CreateBank(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, bank)
}
