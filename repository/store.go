package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/mateuszkochelski/SwiftCodeDb/db/sqlc"
)

func InsertCountryWithValidation(queries *db.Queries, newCountry db.CreateCountryParams) error {
	existingCountry, getError := queries.GetCountry(context.Background(), newCountry.CountryCode)
	if getError == sql.ErrNoRows {
		_, insertionError := queries.CreateCountry(context.Background(), newCountry)
		if insertionError != nil {
			return fmt.Errorf("insertion failed %s", insertionError.Error())
		}
		return nil
	} else if getError != nil {
		return fmt.Errorf("query error %s", getError)
	}

	if existingCountry.CountryName != newCountry.CountryName {
		return errors.New("inconsistency in request, insertion stopped")
	}
	return nil
}

func InsertBankWithValidation(queries *db.Queries, newBank db.CreateBankParams) error {
	_, getError := queries.GetBankBySwiftCodeWithCountry(context.Background(), newBank.SwiftCode)
	if getError == sql.ErrNoRows {
		_, insertionError := queries.CreateBank(context.Background(), newBank)
		if insertionError != nil {
			return fmt.Errorf("insertion failed %s", insertionError.Error())
		}
		return nil
	} else if getError != nil {
		return fmt.Errorf("query error %s", getError)
	}

	return errors.New("there was bank with existing code in database")
}
