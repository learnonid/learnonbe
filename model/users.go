package model

// Struktur untuk pengguna (Users)
type Users struct {
	UserID      uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	FullName    string `gorm:"not null" json:"full_name"`
	Email       string `gorm:"unique;not null" json:"email"`
	PhoneNumber string `gorm:"size:15" json:"phone_number"`
	Password    string `gorm:"not null" json:"password"`
	RoleID      uint   `gorm:"not null" json:"role_id"`
	Status      string `gorm:"default:'biasa'" json:"status"`
	Roles       Roles  `gorm:"foreignKey:RoleID" json:"role"`
}

// Struct untuk Role
type Roles struct {
	RoleID   uint   `gorm:"primaryKey" json:"role_id"`
	RoleName string `gorm:"not null" json:"role_name"`
}
