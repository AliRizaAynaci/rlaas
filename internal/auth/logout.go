package auth

import "github.com/gofiber/fiber/v2"

// GET /logout – deletes cookie, returns 204
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // expire immediately
		HTTPOnly: true,
	})
	return c.SendStatus(fiber.StatusNoContent)
}
