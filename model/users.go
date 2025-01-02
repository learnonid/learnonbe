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
	RoleID      primitive.ObjectID 	`bson:"role_id" json:"role_id"`               // ID role (referensi ke Roles)
	Status      string             	`bson:"status" json:"status"`                 // Status
	CreatedAt	primitive.DateTime	`bson:"created_at,omitempty" json:"created_at"`
}

// Struktur untuk role (Roles)
type Roles struct {
	RoleID   primitive.ObjectID `bson:"_id,omitempty" json:"role_id"`  // ID unik role
	RoleName string             `bson:"role_name" json:"role_name"`   // Nama role
}
