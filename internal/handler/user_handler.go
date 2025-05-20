package handler

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/app"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/util"
	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Service *app.UserService
}

func NewUserHandler(service *app.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	username := util.GetUsernameFromJWT(c)
	if username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	user, err := h.Service.GetProfile(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Success: false,
			Error:   "user not found",
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Data: fiber.Map{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	username := util.GetUsernameFromJWT(c)
	if username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	type req struct {
		Password string `json:"password,omitempty"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	if err := h.Service.UpdateProfile(username, body.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Message: "profile updated",
	})
}
