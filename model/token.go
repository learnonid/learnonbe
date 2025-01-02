package model

import (
	"github.com/golang-jwt/jwt/v4" // Gunakan versi terbaru jwt-go
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struktur untuk token
type JWTClaims struct {
	jwt.RegisteredClaims      // Mengganti StandardClaims dengan RegisteredClaims (di jwt/v4)
	UserID     primitive.ObjectID `json:"user_id"` // ID pengguna dalam format MongoDB ObjectID
	RoleID     primitive.ObjectID `json:"role_id"` // ID peran dalam format MongoDB ObjectID
}
