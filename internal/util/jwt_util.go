package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
	
	// Try to extract from JWT token object first
	if token, ok := user.(*jwt.Token); ok {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return extractFromClaims(claims)
		}
	}
	
	// Fall back to map extraction
	claims, ok := user.(jwt.MapClaims)
	if ok {
		return extractFromClaims(claims)
	}
	
	// Last resort - try as a regular map
	claimsMap, ok := user.(map[string]interface{})
	if ok {
		return extractFromMap(claimsMap)
	}
	
	return nil, false
}

func extractFromClaims(claims jwt.MapClaims) (*UserClaims, bool) {
	var userClaims UserClaims
	
	// Extract username
	if username, exists := claims["username"]; exists {
		if usernameStr, ok := username.(string); ok && usernameStr != "" {
			userClaims.Username = usernameStr
		} else {
			return nil, false
		}
	} else {
		return nil, false
	}
	
	// Extract type
	if userType, exists := claims["type"]; exists {
		if typeStr, ok := userType.(string); ok {
			userClaims.Type = typeStr
		} else {
			userClaims.Type = "user"
		}
	} else {
		userClaims.Type = "user"
	}
	
	// Extract user_id
	if userID, exists := claims["user_id"]; exists {
		switch v := userID.(type) {
		case float64:
			userClaims.UserID = uint(v)
		case int:
			userClaims.UserID = uint(v)
		case uint:
			userClaims.UserID = v
		default:
			// If we can't extract user_id properly, return false
			return nil, false
		}
	}
	
	return &userClaims, true
}

func extractFromMap(claims map[string]interface{}) (*UserClaims, bool) {
	var userClaims UserClaims
	
	// Extract username
	if username, exists := claims["username"]; exists {
		if usernameStr, ok := username.(string); ok && usernameStr != "" {
			userClaims.Username = usernameStr
		} else {
			return nil, false
		}
	} else {
		return nil, false
	}
	
	// Extract type
	if userType, exists := claims["type"]; exists {
		if typeStr, ok := userType.(string); ok {
			userClaims.Type = typeStr
		} else {
			userClaims.Type = "user"
		}
	} else {
		userClaims.Type = "user"
	}
	
	// Extract user_id
	if userID, exists := claims["user_id"]; exists {
		switch v := userID.(type) {
		case float64:
			userClaims.UserID = uint(v)
		case int:
			userClaims.UserID = uint(v)
		case uint:
			userClaims.UserID = v
		default:
			// If we can't extract user_id properly, return false
			return nil, false
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
