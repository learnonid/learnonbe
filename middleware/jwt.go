package middleware

import (
	"learnonbe/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Periksa apakah token ada di header Authorization
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "No token provided")
		}

		// Parse token dan validasi
		token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// Ekstrak claims dari token
		claims, ok := token.Claims.(*model.JWTClaims)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse claims")
		}

		// Simpan claims dalam context untuk digunakan di handler berikutnya
		c.Locals("claims", claims)

		return c.Next() // Lanjutkan ke handler berikutnya
	}
}