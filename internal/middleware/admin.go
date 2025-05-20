package middleware

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/util"
	"github.com/gofiber/fiber/v2"
)

func AdminOnly(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "admin access required",
		})
	}

	return c.Next()
}
