package handler

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/app"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/util"
	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/request"
	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Service *app.AuthService
}

func NewAuthHandler(service *app.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "username and password required",
		})
	}
	if err := h.Service.Register(body.Username, body.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(response.Response{
		Success: true,
		Message: "user registered",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	token, err := h.Service.Login(body.Username, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Data:    response.LoginData{Token: token},
	})
}

func (h *AuthHandler) RegisterAdmin(c *fiber.Ctx) error {
	var body request.AdminRegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "username and password required",
		})
	}
	if err := h.Service.RegisterWithType(body.Username, body.Password, "admin"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(response.Response{
		Success: true,
		Message: "admin registered",
	})
}

func (h *AuthHandler) LoginAdmin(c *fiber.Ctx) error {
	var body request.AdminLoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	token, err := h.Service.LoginWithType(body.Username, body.Password, "admin")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Data:    response.LoginData{Token: token},
	})
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userName := util.GetUsernameFromJWT(c)
	if userName == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	user, err := h.Service.GetProfile(userName)
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
			"type":     user.Type,
		},
	})
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userName := util.GetUsernameFromJWT(c)
	if userName == "" {
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
	if body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "password required",
		})
	}
	if err := h.Service.UpdateProfile(userName, body.Password); err != nil {
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
