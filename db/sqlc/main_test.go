package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver         = "postgres"
	dbSource         = "postgresql://test:test@postgresTestDB:5432/testdb?sslmode=disable"
	validCountryCode = "PL"
	validCountryname = "POLAND"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, _ := sql.Open(dbDriver, dbSource)
	defer conn.Close()
	err := conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
