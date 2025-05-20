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
	
	username, exists := claims["username"]
	if !exists {
		return nil, false
	}
	usernameStr, ok := username.(string)
	if !ok || usernameStr == "" {
		return nil, false
	}
	userClaims.Username = usernameStr
	
	userType, exists := claims["type"]
	if exists {
		if typeStr, ok := userType.(string); ok {
			userClaims.Type = typeStr
		} else {
			userClaims.Type = "user" 
		}
	} else {
		userClaims.Type = "user" 
	}
	
	if userID, exists := claims["user_id"]; exists {
		if idFloat, ok := userID.(float64); ok {
			userClaims.UserID = uint(idFloat)
		}
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

func GetUserIDFromJWT(c *fiber.Ctx) uint {
	user, ok := ExtractUser(c)
	if !ok {
		return 0
	}
	return user.UserID
}
