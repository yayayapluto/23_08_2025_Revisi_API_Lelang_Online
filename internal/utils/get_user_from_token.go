package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"log"
)

func GetUserFromToken(c *fiber.Ctx) (*entities.User, error) {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok && token == nil {
		return nil, errors.New("no token found")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	var organizerId *uint
	if organizerIdToken, ok := claims["organizer_id"].(float64); ok {
		id := uint(organizerIdToken)
		if id > 0 {
			organizerId = &id
		}
	}

	log.Println(claims)

	user := &entities.User{
		ID:          uint(claims["user_id"].(float64)),
		Username:    claims["username"].(string),
		Email:       claims["email"].(string),
		Role:        claims["role"].(string),
		OrganizerID: organizerId,
	}

	return user, nil
}
