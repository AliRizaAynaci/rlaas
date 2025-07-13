package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET /logout – deletes cookie, returns 204
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Domain:   ".rlaas.tech",   // ← başlangıçta set ettiğiniz domain
		Path:     "/",             // ← başlangıçta da "/"
		Expires:  time.Unix(0, 0), // veyahut MaxAge: -1
		HTTPOnly: true,
		SameSite: "None", // ← başlangıçta da None
		Secure:   true,   // ← HTTPS zorunlu
		MaxAge:   -1,
	})
	return c.SendStatus(fiber.StatusNoContent)
}
