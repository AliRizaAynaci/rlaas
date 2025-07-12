package rule

import "time"

type Rule struct {
	ID            uint      `json:"id"            gorm:"primaryKey"`
	ProjectID     uint      `json:"project_id"    gorm:"index"`
	Endpoint      string    `json:"endpoint"`
	Strategy      string    `json:"strategy"` // token_bucket | sliding_window | …
	KeyBy         string    `json:"key_by"`   // api_key | ip | user_id
	LimitCount    int       `json:"limit_count"`
	WindowSeconds int       `json:"window_seconds"`
	CreatedAt     time.Time `json:"created_at"`
}
