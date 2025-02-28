package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
	store "github.com/mateuszkochelski/SwiftCodeDb/repository"
)

var queries *db.Queries

const (
	csvPath           = "swift_codes.csv"
	swiftCodeLenght   = 11
	countryCodeLenght = 2
	numberOfColumns   = 8
	dbDriver          = "postgres"
	dbSource          = "postgresql://root:password@postgresDB:5432/swift_codes?sslmode=disable"
)

func getBankType(swiftCode string) db.BankType {
	if strings.HasSuffix(swiftCode, "XXX") {
		return db.BankTypeHeadquarter
	}
	return db.BankTypeBranch
}

func validateData(countryCode, swiftCode, bankName, countryName string) error {

	var errs string
	if len(countryCode) != countryCodeLenght {
		errs += "country code must be lenght of 2,"
	}
	if strings.ToUpper(countryCode) != countryCode {
		errs += "country code must be uppercase,"
	}
	if len(swiftCode) != swiftCodeLenght {
		errs += "swift code must be lenght of 11,"
	}
	if strings.ToUpper(countryName) != countryName {
		errs += "country names must be uppercase,"
	}
	if len(bankName) == 0 {
		errs += "bank name must be not null"
	}
	if errs != "" {
		return errors.New(strings.TrimSuffix(errs, ","))
	}
	return nil
}

func getDataFromRecord(record []string) (db.CreateBankParams, db.CreateCountryParams, error) {
	if len(record) != numberOfColumns {
		return db.CreateBankParams{}, db.CreateCountryParams{}, errors.New("invalid numbers of column in record")
	}
	countryCode := record[0]
	swiftCode := record[1]
	bankName := record[3]
	bankAddress := record[4]
	countryName := record[6]
	err := validateData(countryCode, swiftCode, bankName, countryName)
	if err != nil {
		return db.CreateBankParams{}, db.CreateCountryParams{}, fmt.Errorf("invalid data: %s", err.Error())
	}
	bankType := getBankType(swiftCode)
	bank := db.CreateBankParams{
		SwiftCode:   swiftCode,
		BankName:    bankName,
		BankAddress: sql.NullString{String: bankAddress, Valid: len(bankAddress) != 0},
		CountryCode: countryCode,
		BankType:    bankType,
	}
	country := db.CreateCountryParams{CountryCode: countryCode, CountryName: countryName}

	return bank, country, nil
}

func main() {
	conn, _ := sql.Open(dbDriver, dbSource)
	err := conn.Ping()
	if err != nil {
		log.Fatal("Error during connection with database")
		return
	}
	defer conn.Close()
	queries = db.New(conn)

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("Error during opening csv")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	line := 0
	for {
		record, err := reader.Read()
		if line == 0 {
			line++
			continue
		}
		if err != nil {
			break
		}
		bank, country, err := getDataFromRecord(record)
		if err != nil {
			fmt.Printf("Invalid data at row %d : %s", line, err.Error())
		}

		err = store.InsertCountryWithValidation(queries, country)
		if err != nil {
			fmt.Printf("Invalid data at row %d : %s", line, err.Error())
		}
		err = store.InsertBankWithValidation(queries, bank)
		if err != nil {
			fmt.Printf("Invalid data at row %d : %s", line, err.Error())
		}
		line++
	}
}
