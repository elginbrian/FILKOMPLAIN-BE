package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTProtected(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(secret),
		ContextKey:   "user",
		ErrorHandler: jwtError,
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"error":   "Unauthorized",
	})
}
