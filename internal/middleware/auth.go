package middleware

import (
	"os"
	"strings"

	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret fetches the secret key
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set!")
	}
	return []byte(secret)
}

// Protected parses the JWT from Authorization header and injects claims into fiber context
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Missing or invalid token. Must be format 'Bearer <token>'",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validate algorithm
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid token claims structure",
			})
		}

		// Inject to locals
		c.Locals("user_id", claims["user_id"])
		c.Locals("is_admin", claims["is_admin"])

		return c.Next()
	}
}

// AdminOnly ensures the user has is_admin set to true. MUST be used after Protected()
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(response.WebResponse{
				Code:    fiber.StatusForbidden,
				Message: "Forbidden: Admin access required",
			})
		}
		return c.Next()
	}
}
