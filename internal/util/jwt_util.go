package util

import (
	"github.com/gofiber/fiber/v2"
)

type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

func ExtractUser(c *fiber.Ctx) (*UserClaims, bool) {
	user := c.Locals("user")
	if user == nil {
		return nil, false
	}
	
	claims, ok := user.(map[string]interface{})
	if !ok {
		return nil, false
	}
	
	var userClaims UserClaims
	
	if username, exists := claims["username"].(string); exists {
		userClaims.Username = username
	} else {
		return nil, false
	}
	
	if userType, exists := claims["type"].(string); exists {
		userClaims.Type = userType
	} else {
		userClaims.Type = "user"
	}
	
	if userID, exists := claims["user_id"].(float64); exists {
		userClaims.UserID = uint(userID)
	}
	
	return &userClaims, true
}

func IsAdmin(c *fiber.Ctx) bool {
	user, ok := ExtractUser(c)
	if !ok {
		return false
	}
	return user.Type == "admin"
}

func GetUsernameFromJWT(c *fiber.Ctx) string {
	user, ok := ExtractUser(c)
	if !ok {
		return ""
	}
	return user.Username
}
