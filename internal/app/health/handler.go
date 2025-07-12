package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Handler holds slow-changing deps we want to ping.
type Handler struct{ db *gorm.DB }

// New returns a health handler with injected DB.
func New(db *gorm.DB) *Handler { return &Handler{db: db} }

/* ---------- /healthz : liveness ---------- */
func (h *Handler) Liveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

/* ---------- /readyz : readiness ---------- */
func (h *Handler) Readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()

	// simple SELECT 1 ping
	if err := h.db.WithContext(ctx).Raw("SELECT 1").Error; err != nil {
		return fiber.ErrServiceUnavailable
	}
	return c.JSON(fiber.Map{"status": "ready"})
}
