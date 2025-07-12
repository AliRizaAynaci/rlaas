package rule

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ svc *Service }

func NewHandler(s *Service) *Handler { return &Handler{s} }

/* helpers */
func pid(c *fiber.Ctx) uint { id, _ := strconv.Atoi(c.Params("pid")); return uint(id) }
func rid(c *fiber.Ctx) uint { id, _ := strconv.Atoi(c.Params("rid")); return uint(id) }

/* GET /projects/:pid/rules */
func (h *Handler) List(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(uint)
	out, err := h.svc.List(uid, pid(c))
	if err != nil {
		return fiber.ErrForbidden
	}
	return c.JSON(out)
}

/* POST /projects/:pid/rules */
func (h *Handler) Create(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(uint)
	var in Rule
	if err := c.BodyParser(&in); err != nil {
		return fiber.ErrBadRequest
	}
	r, err := h.svc.Add(uid, pid(c), &in)
	if err != nil {
		return fiber.ErrForbidden
	}
	return c.Status(fiber.StatusCreated).JSON(r)
}

/* PUT /projects/:pid/rules/:rid */
func (h *Handler) Update(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(uint)
	var in Rule
	if err := c.BodyParser(&in); err != nil {
		return fiber.ErrBadRequest
	}
	in.ID = rid(c)
	in.ProjectID = pid(c)
	if err := h.svc.Update(uid, &in); err != nil {
		return fiber.ErrForbidden
	}
	return c.SendStatus(fiber.StatusOK)
}

/* DELETE /projects/:pid/rules/:rid */
func (h *Handler) Delete(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(uint)
	if err := h.svc.Delete(uid, pid(c), rid(c)); err != nil {
		return fiber.ErrForbidden
	}
	return c.SendStatus(fiber.StatusNoContent)
}
