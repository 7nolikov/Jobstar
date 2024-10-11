package db

import (
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() {
    connStr := os.Getenv("DB_CONNECTION_STRING")
    if connStr == "" {
        connStr = "postgres://user:password@localhost:5432/yourdb?sslmode=disable"
    }

    var err error
    DB, err = sqlx.Connect("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    fmt.Println("Database connection established")
}
