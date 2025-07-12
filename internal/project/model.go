package project

import (
	"time"

	"github.com/AliRizaAynaci/rlaas/internal/rule"
)

type Project struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    uint        `json:"user_id" gorm:"index"`
	Name      string      `json:"name"`
	APIKey    string      `json:"api_key" gorm:"uniqueIndex"`
	CreatedAt time.Time   `json:"created_at"`
	Rules     []rule.Rule `json:"rules" gorm:"constraint:OnDelete:CASCADE"`
}
