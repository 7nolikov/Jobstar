package models

import "time"

type Vacancy struct {
    ID          int       `db:"id"`
    Title       string    `db:"title"`
    Description string    `db:"description"`
    Location    string    `db:"location"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}
