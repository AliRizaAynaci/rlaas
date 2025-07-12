package user

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GoogleID  string    `json:"google_id" gorm:"uniqueIndex"`
	Email     string    `json:"email"     gorm:"uniqueIndex"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}
