package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:8080/swift_codes?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, _ := sql.Open(dbDriver, dbSource)
	err := conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
