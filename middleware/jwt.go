package middleware

import (
	"learnonbe/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"fmt"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "No token provided")
		}

		// Parse token dan klaim
		token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// Ambil klaim
		claims, ok := token.Claims.(*model.JWTClaims)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse claims")
		}

		fmt.Printf("Parsed Claims: UserID=%d, RoleID=%d\n", claims.UserID, claims.RoleID)

		// Simpan klaim ke context
		c.Locals("claims", claims)
		return c.Next()
	}
}
