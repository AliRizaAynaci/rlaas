package user

import "github.com/gofiber/fiber/v2"

type Handler struct{ svc *Service }

func NewHandler(s *Service) *Handler { return &Handler{s} }

func (h *Handler) Me(c *fiber.Ctx) error {
	uidAny := c.Locals("user_id")
	if uidAny == nil {
		return fiber.ErrUnauthorized
	}
	user, err := h.svc.GetByID(uidAny.(uint))
	if err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}
