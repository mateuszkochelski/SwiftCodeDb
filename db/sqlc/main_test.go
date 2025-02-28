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
	dbSource         = "postgresql://root:password@localhost:5432/swift_codes?sslmode=disable"
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
