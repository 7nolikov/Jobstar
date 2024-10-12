package db

import (
    "log"
    "os"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
    connStr := os.Getenv("DB_CONNECTION_STRING")
    if connStr == "" {
        log.Fatal("DB_CONNECTION_STRING is not set")
    }

    log.Printf("Connecting to database with: %s\n", connStr) // Debugging line

    var err error
    DB, err = sqlx.Connect("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("Database connection established")
}
