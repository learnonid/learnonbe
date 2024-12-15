package model

import "github.com/dgrijalva/jwt-go"

// Struktur untuk token
type JWTClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}