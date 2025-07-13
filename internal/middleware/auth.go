package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Auth validates the JWT found in cookie or Authorization header
func Auth() fiber.Handler {
	secret := []byte(getenv("JWT_SECRET", "super-secret-change-me"))

	return func(c *fiber.Ctx) error {
		var tokenStr string

		// 1) Try session cookie first
		if cookie := c.Cookies("session_token"); cookie != "" {
			tokenStr = cookie
		}
		// 2) Fallback to Authorization: Bearer <token> header
		if tokenStr == "" {
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenStr = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		if tokenStr == "" {
			log.Fatal("TOKENSTR IS EMPTY")
			return fiber.ErrUnauthorized
		}

		// Parse and validate the token
		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !tok.Valid {
			log.Fatal("JWT PARSE FAILED: ", err)
			return fiber.ErrUnauthorized
		}

		claims := tok.Claims.(jwt.MapClaims)
		uid := uint(claims["user_id"].(float64))
		// Pass user_id to the next handlers
		c.Locals("user_id", uid)

		return c.Next()
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
