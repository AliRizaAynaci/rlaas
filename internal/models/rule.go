package models

import "time"

// RateLimitRule ties a specific endpoint under a Project to its rate-limit config.
type RateLimitRule struct {
	ID            int       `db:"id"`
	ProjectID     int       `db:"project_id"` // FK to projects.id
	Endpoint      string    `db:"endpoint"`   // e.g. "/api/v1/check"
	Strategy      string    `db:"strategy"`   // e.g. "sliding_window"
	KeyBy         string    `db:"key_by"`     // e.g. "api_key" or "ip"
	LimitCount    int       `db:"limit_count"`
	WindowSeconds int       `db:"window_seconds"`
	CreatedAt     time.Time `db:"created_at"`
}
