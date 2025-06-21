package models

import "time"

type RateLimitRule struct {
	ID            int
	ProjectID     int
	Endpoint      string
	Strategy      string
	KeyBy         string
	LimitCount    int
	WindowSeconds int
	CreatedAt     time.Time
}
