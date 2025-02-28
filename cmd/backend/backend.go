package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	handlers "github.com/mateuszkochelski/SwiftCodeDb/handlers"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@postgresDB:5432/swift_codes?sslmode=disable"
)

func main() {
	conn, _ := sql.Open(dbDriver, dbSource)
	defer conn.Close()
	err := conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	bankHandler := handlers.NewBankHandler(conn)

	http.HandleFunc("/v1/swift-codes/", bankHandler.HandleSwiftCodes)
	http.HandleFunc("/v1/swift-codes/country/", bankHandler.GetBanksByContryCode)
	http.HandleFunc("/v1/swift-codes", bankHandler.CreateBank)
	http.ListenAndServe(":8080", nil)

}
