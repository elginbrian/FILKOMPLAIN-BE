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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "username, email, and password required",
		})
	}
	if err := h.Service.Register(body.Username, body.Email, body.Password); err != nil {
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
		Email    string `json:"email"` 
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	token, err := h.Service.Login(body.Email, body.Password) 
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
	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "username, email, and password required",
		})
	}
	if err := h.Service.RegisterWithType(body.Username, body.Email, body.Password, "admin"); err != nil {
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
	token, err := h.Service.LoginWithType(body.Email, body.Password, "admin") // Changed to use Email
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
	userID := util.GetUserIDFromJWT(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	user, err := h.Service.GetProfileByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Success: false,
			Error:   "user not found",
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Data: fiber.Map{
			"id":            user.ID,
			"username":      user.Username,
			"email":         user.Email,
			"type":          user.Type,
			"nim":           user.NIM,
			"program_studi": user.ProgramStudi,
			"phone_number":  user.PhoneNumber,
		},
	})
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := util.GetUserIDFromJWT(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	
	var body request.UpdateProfileRequest
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
	
	if err := h.Service.UpdateProfileFull(userID, body.Password, body.NIM, body.ProgramStudi, body.PhoneNumber); err != nil {
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

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	username := util.GetUsernameFromJWT(c)
	if username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	
	user, err := h.Service.GetProfile(username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "user not found",
		})
	}
	
	token, err := h.Service.GenerateTokenForUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   "could not refresh token",
		})
	}
	
	return c.JSON(response.Response{
		Success: true,
		Data:    response.LoginData{Token: token},
	})
}
