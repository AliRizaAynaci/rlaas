package models

import "time"

// User represents a registered user (via Google OAuth).
type User struct {
	ID        int       `db:"id"`
	GoogleID  string    `db:"google_id"` // unique Google account identifier
	Email     string    `db:"email"`     // unique user email
	CreatedAt time.Time `db:"created_at"`
}
