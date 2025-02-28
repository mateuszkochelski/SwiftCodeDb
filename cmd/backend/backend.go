package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	handlers "github.com/mateuszkochelski/SwiftCodeDb/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DB_SOURCE")

	conn, _ := sql.Open(dbDriver, dbSource)
	defer conn.Close()
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	bankHandler := handlers.NewBankHandler(conn)

	http.HandleFunc("/v1/swift-codes/", bankHandler.HandleSwiftCodes)
	http.HandleFunc("/v1/swift-codes/country/", bankHandler.GetBanksByContryCode)
	http.HandleFunc("/v1/swift-codes", bankHandler.CreateBank)
	http.ListenAndServe(":8080", nil)

}
