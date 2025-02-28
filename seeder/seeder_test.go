package main

import (
	"database/sql"
	"testing"

	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	"github.com/mateuszkochelski/SwiftCodeDb/util"
	"github.com/stretchr/testify/require"
)

func Test_get_bank_type(t *testing.T) {

	tests := []struct {
		name      string
		swiftCode string
		bankType  db.BankType
	}{
		{
			name:      "returning_headquarter_given_swift_code_ended_with_XXX",
			swiftCode: "12312312XXX",
			bankType:  db.BankTypeHeadquarter,
		},
		{
			name:      "returning_headquarter_given_swift_code_with_X_letters_only_longer_than_2",
			swiftCode: "XXXXXXXXXXXXXXXXXXXXXXXXX",
			bankType:  db.BankTypeHeadquarter,
		},
		{
			name:      "returning_branch_given_swift_code_not_long_enough",
			swiftCode: "XX",
			bankType:  db.BankTypeBranch,
		},
		{
			name:      "returning_branch_given_swift_code_with_XXX_in_middle",
			swiftCode: "XXABCXXXABFCD",
			bankType:  db.BankTypeBranch,
		},
		{
			name:      "returning_branch_given_swift_code_without_X_letter",
			swiftCode: "123123123123123",
			bankType:  db.BankTypeBranch,
		},
		{
			name:      "returning_branch_given_swift_code_empty",
			swiftCode: "",
			bankType:  db.BankTypeBranch,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getBankType(test.swiftCode)
			require.Equal(t, test.bankType, result)
		})
	}
}

func Test_validate_data_looking_for_error_tests(t *testing.T) {
	tests := []struct {
		name         string
		countryCode  string
		swiftCode    string
		bankName     string
		countryName  string
		errorMessage string
	}{
		{
			name:         "returning_wrong_country_code_lenght_error1",
			countryCode:  "A",
			errorMessage: "country code must be lenght of 2",
		},
		{
			name:         "returning_wrong_country_code_lenght_error2",
			countryCode:  "AAA",
			errorMessage: "country code must be lenght of 2",
		},
		{
			name:         "returning_wrong_country_code_lenght_error3",
			countryCode:  "AAAAAA",
			errorMessage: "country code must be lenght of 2",
		},
		{
			name:         "returning_country_code_must_be_uppercase_error",
			countryCode:  "AAAAAAa",
			errorMessage: "country code must be uppercase",
		},
		{
			name:         "returning_swift_code_wrong_lenght_error",
			swiftCode:    "ABD",
			errorMessage: "swift code must be lenght of 11",
		},
		{
			name:         "returning_country_names_must_up_be_uppercase_error",
			countryName:  "poland",
			errorMessage: "country names must be uppercase",
		},
		{
			name:         "return_bank_name_must_be_not_null_error",
			errorMessage: "bank name must be not null",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateData(
				test.countryCode,
				test.swiftCode,
				test.bankName,
				test.countryName,
			)
			require.Contains(t, err.Error(), test.errorMessage)
		})
	}
}

func Test_validate_data_looking_for_no_specific_errors(t *testing.T) {
	tests := []struct {
		name         string
		countryCode  string
		swiftCode    string
		bankName     string
		countryName  string
		errorMessage string
	}{
		{
			name:         "not_returning_wrong_country_code_lenght_error",
			countryCode:  "AA",
			errorMessage: "Country code must be lenght of 2",
		},
		{
			name:         "not_returning_country_code_must_be_uppercase_error",
			countryCode:  "PL",
			errorMessage: "Country code must be uppercase",
		},
		{
			name:         "not_returning_swift_code_wrong_lenght_error",
			swiftCode:    "ABCDEFGH123",
			errorMessage: "Swift code must be lenght of 11",
		},
		{
			name:         "not_returning_country_names_must_up_be_uppercase_error",
			countryName:  "POLAND",
			errorMessage: "Country names must be uppercase",
		},
		{
			name:         "not_return_bank_name_must_be_not_null_error",
			bankName:     "Pekao",
			errorMessage: "Bank name must be not null",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateData(
				test.countryCode,
				test.swiftCode,
				test.bankName,
				test.countryName,
			)
			require.NotContains(t, err.Error(), test.errorMessage)
		})
	}
}

func Test_validate_data_looking_for_no_errors_at_all(t *testing.T) {
	test := struct {
		name        string
		countryCode string
		swiftCode   string
		bankName    string
		countryName string
	}{
		name:        "returning_no_errors_given_valid_input",
		countryCode: "AL",
		swiftCode:   "12345678XXX",
		bankName:    "pekao",
		countryName: "ALBANIA",
	}

	t.Run(test.name, func(t *testing.T) {
		err := validateData(
			test.countryCode,
			test.swiftCode,
			test.bankName,
			test.countryName,
		)
		require.NoError(t, err)
	})
}

func Test_validate_data_looking_for_errors_at_all(t *testing.T) {
	test := struct {
		name        string
		countryCode string
		swiftCode   string
		bankName    string
		countryName string
	}{
		name:        "returning_error_given_invalid_input",
		countryCode: "PL",
		swiftCode:   "12345678XXXX",
		bankName:    "pekao",
		countryName: "ALBANIA",
	}

	t.Run(test.name, func(t *testing.T) {
		err := validateData(
			test.countryCode,
			test.swiftCode,
			test.bankName,
			test.countryName,
		)
		require.Error(t, err)
	})
}

func Test_get_data_from_record_returns_wrong_number_of_column_error(t *testing.T) {
	var columns []string
	const invalid_columns_error_message = "invalid numbers of column in record"
	for i := 0; i < 7; i++ {
		columns = append(columns, util.RandomString(10))
	}
	_, _, err := getDataFromRecord(columns)
	require.EqualError(t, err, invalid_columns_error_message)

	columns = append(columns, util.RandomString(2))
	_, _, err = getDataFromRecord(columns)
	require.NotContains(t, err.Error(), invalid_columns_error_message)

	columns = append(columns, util.RandomString(2))
	_, _, err = getDataFromRecord(columns)
	require.EqualError(t, err, invalid_columns_error_message)
}

func Test_get_data_from_record_returning_country_and_bank_given_valid_data(t *testing.T) {
	testRecord := []string{
		"AL", "12345678XXX", "BIC11", "Bank", "", "cravow", "ALBANIA", "Pacific",
	}
	bank, country, err := getDataFromRecord(testRecord)
	require.NoError(t, err)
	require.Equal(t, bank.CountryCode, "AL")
	require.Equal(t, bank.SwiftCode, "12345678XXX")
	require.Equal(t, bank.BankAddress, sql.NullString{String: "", Valid: false})
	require.Equal(t, bank.BankName, "Bank")
	require.Equal(t, country.CountryCode, "AL")
	require.Equal(t, country.CountryName, "ALBANIA")
}
