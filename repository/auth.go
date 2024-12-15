package repository

import (
	"learnonbe/model"
	"learnonbe/utils"

	"time"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateRole(db *gorm.DB, role *model.Roles) error {
	// Buat jika RoleID 1 (admin) jika RoleID 2 (user)
	if role.RoleID == 1 {
		// Buat user admin
		role.RoleName = "admin"
	} else if role.RoleID == 2 {
		// Buat user
		role.RoleName = "user"
	} 

	// Buat Role
	err := db.Create(&role).Error
	if err != nil {
		return fmt.Errorf("Gagal membuat role: %v", err)
	}
	
	return nil
}

func CreateUser(db *gorm.DB, user *model.Users) error {
	// Generate UserID secara random
	user.UserID = utils.GenerateRandomID(1, 10000)

	// Validasi RoleID jika tidak diisi maka set default sebagai 2 (user)
	if user.RoleID == 0 {
		user.RoleID = 2 // user
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Ubah jika user menginput nomor telepon dengan 08 menjadi +628
	user.PhoneNumber = ChangePhoneNumber(user.PhoneNumber)

	// Validasi phone number
	if err := ValidatePhoneNumber(user.PhoneNumber); err != nil {
		return err
	}

	// Validasi email
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	// Set status
	user.Status = "biasa"


	// Create user
	err = db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func GenerateToken(UserID uint) (string, error) {
	// Claim
	claims := jwt.MapClaims{
		"UserID": UserID,
		"exp": time.Now().Add(time.Hour * 12).Unix(), // Token berlaku selama 12 jam
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}