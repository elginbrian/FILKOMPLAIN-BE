package handler

import (
	"mime/multipart"

	"github.com/elginbrian/FILKOMPLAIN-BE/internal/app"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/util"
	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/request"
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
			"id":               user.ID,
			"username":         user.Username,
			"email":            user.Email,
			"type":             user.Type,
			"nim":              user.NIM,
			"program_studi":    user.ProgramStudi,
			"phone_number":     user.PhoneNumber,
			"profile_image_url": user.ProfileImageURL,
		},
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := util.GetUserIDFromJWT(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Success: false,
			Error:   "unauthorized",
		})
	}
	
	// Check if this is a multipart form request (with file upload)
	if form, err := c.MultipartForm(); err == nil {
		var body request.UpdateProfileRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Response{
				Success: false,
				Error:   "invalid request",
			})
		}
		
		var profileImage *multipart.FileHeader
		if files := form.File["profile_image"]; len(files) > 0 {
			profileImage = files[0]
		}
		
		if err := h.Service.UpdateProfileWithImage(userID, body.Password, body.Username, body.NIM, body.ProgramStudi, body.PhoneNumber, profileImage); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Success: false,
				Error:   err.Error(),
			})
		}
		
		return c.JSON(response.Response{
			Success: true,
			Message: "profile updated with image",
		})
	}
	
	// Handle regular JSON request (no file upload)
	var body request.UpdateProfileRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid request",
		})
	}
	
	if err := h.Service.UpdateProfileFull(userID, body.Password, body.Username, body.NIM, body.ProgramStudi, body.PhoneNumber); err != nil {
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
