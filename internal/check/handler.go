package check

import (
	"github.com/gofiber/fiber/v2"

	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"github.com/AliRizaAynaci/rlaas/internal/service"
)

type Handler struct{ svc *service.RateConfigService }

func NewHandler(s *service.RateConfigService) *Handler { return &Handler{svc: s} }

func (h *Handler) Handle(c *fiber.Ctx) error {
	var req struct {
		APIKey   string `json:"api_key"`
		Endpoint string `json:"endpoint"`
		Key      string `json:"key"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	cfg, err := h.svc.Get(req.APIKey, req.Endpoint)
	switch err {
	case service.ErrProjectNotFound:
		return fiber.ErrUnauthorized
	case service.ErrEndpointNotOwned:
		return fiber.ErrForbidden
	case nil:
	default:
		return fiber.ErrInternalServerError
	}

	lim, _ := limiter.GetLimiterForKey(req.APIKey, req.Endpoint, req.Key, cfg)
	allowed, _ := lim.Allow(req.Key)
	if !allowed {
		return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
	}

	return c.JSON(fiber.Map{"allowed": true})
}
