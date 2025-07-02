package models

import "time"

// Project represents a namespace for rate-limit rules, owned by a User.
type Project struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"` // FK to users.id
	Name      string    `db:"name"`    // unique per user
	ApiKey    string    `db:"api_key"` // generated 256-bit key
	CreatedAt time.Time `db:"created_at"`
}
