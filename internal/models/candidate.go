package models

import "time"

type Candidate struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	Resume    string `db:"resume" json:"resume"`
	State     string    `db:"state" json:"state"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
