package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struktur untuk pengguna (Users)
type Users struct {
	UserID      primitive.ObjectID 	`bson:"_id,omitempty" json:"user_id"`          // ID unik pengguna
	FullName    string             	`bson:"full_name" json:"full_name"`           // Nama lengkap pengguna
	Email       string             	`bson:"email" json:"email"`                   // Email pengguna
	PhoneNumber string             	`bson:"phone_number,omitempty" json:"phone_number"` // Nomor telepon pengguna
	Password    string             	`bson:"password" json:"password"`             // Kata sandi (hashed)
	RoleID      int                 `bson:"role_id" json:"role_id"`               // ID role (1 untuk admin, 2 untuk customer)
	CreatedAt	primitive.DateTime	`bson:"created_at,omitempty" json:"created_at"`
}

// Struktur untuk role (Roles)
type Roles struct {
	RoleID   int    `bson:"role_id" json:"role_id"`  // ID unik role (1 untuk admin, 2 untuk customer)
	RoleName string `bson:"role_name" json:"role_name"`   // Nama role
}

// Definisikan role yang valid
const (
	RoleAdmin   = 1
	RoleCustomer = 2
)
