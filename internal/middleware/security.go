package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SecurityHeaders applies HTTP security headers using Helmet
func SecurityHeaders() fiber.Handler {
	return helmet.New()
}

// Cors applies Cross-Origin Resource Sharing rules
func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Idempotency-Key",
	})
}

// GlobalRateLimiter applies a DDoS protection rate limit globally
func GlobalRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    fiber.StatusTooManyRequests,
				"message": "Too many requests. Please try again later.",
			})
		},
	})
}

// AuthRateLimiter applies a strict rate limit for authentication endpoints to prevent brute-force
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    fiber.StatusTooManyRequests,
				"message": "Too many login attempts. Locked for 1 minute.",
			})
		},
	})
}
