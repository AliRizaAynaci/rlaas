package project

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ svc *Service }

func NewHandler(s *Service) *Handler { return &Handler{s} }

// POST /projects  { "project_name": "My API" }
func (h *Handler) Create(c *fiber.Ctx) error {
	var req struct {
		ProjectName string `json:"project_name"`
	}
	if err := c.BodyParser(&req); err != nil || req.ProjectName == "" {
		return fiber.ErrBadRequest
	}

	apikey, err := GenerateAPIKey()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	uid := c.Locals("user_id").(uint)
	p, err := h.svc.Create(uid, req.ProjectName, apikey)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

// GET /projects
func (h *Handler) List(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(uint)
	list, err := h.svc.List(uid)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(list)
}

// DELETE /projects/:pid
func (h *Handler) Delete(c *fiber.Ctx) error {
	pid, _ := strconv.Atoi(c.Params("pid"))
	uid := c.Locals("user_id").(uint)

	if err := h.svc.Delete(uid, uint(pid)); err != nil {
		return fiber.ErrForbidden // ya sahibi değil ya proje yok
	}
	return c.SendStatus(fiber.StatusNoContent)
}
