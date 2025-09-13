package middleware

import (
	"github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"log"
	"strings"
)

func Protected(roles ...string) fiber.Handler {
	env, err := utils.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(env.JWT_SECRET)},
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			if len(roles) == 0 {
				return c.Next()
			}

			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)

			roleClaim, ok := claims["role"].(string)
			if !ok {
				return presenters.ErrorResponse(c, fiber.StatusForbidden, "role not found in token", nil)
			}

			for _, r := range roles {
				if strings.EqualFold(r, roleClaim) {
					return c.Next()
				}
			}

			return presenters.ErrorResponse(c, fiber.StatusForbidden, "you don't have permission for this route", nil)
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, "missing or malformed token", err)
	}
	return presenters.ErrorResponse(c, fiber.StatusUnauthorized, "you don't have access", err)
}
