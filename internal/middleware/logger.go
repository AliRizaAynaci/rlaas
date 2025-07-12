package middleware

import (
	"time"

	"github.com/AliRizaAynaci/rlaas/internal/logging"
	"github.com/gofiber/fiber/v2"
)

// RequestLogger returns a Fiber middleware that logs every req/res.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		uid, _ := c.Locals("user_id").(uint)

		switch {
		case status >= 500:
			logging.L.Error("request completed",
				"method", method,
				"path", path,
				"status", status,
				"latency_ms", latency.Milliseconds(),
				"ip", ip,
				"user_id", uid,
				"err", err,
			)
		case status >= 400:
			logging.L.Warn("request completed",
				"method", method,
				"path", path,
				"status", status,
				"latency_ms", latency.Milliseconds(),
				"ip", ip,
				"user_id", uid,
			)
		default:
			logging.L.Info("request completed",
				"method", method,
				"path", path,
				"status", status,
				"latency_ms", latency.Milliseconds(),
				"ip", ip,
				"user_id", uid,
			)
		}

		return err
	}
}
