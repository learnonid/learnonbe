package repository

import (
	"learnonbe/model"
	"learnonbe/utils"

	"time"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	// "golang.org/x/crypto/bcrypt"
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
		return fmt.Errorf("gagal membuat role: %v", err)
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

    // Simpan user tanpa hashing ulang
    err := db.Create(&user).Error
    if err != nil {
        return err
    }

    return nil
}

func GenerateToken(userID, roleID uint) (string, error) {
	// Buat klaim menggunakan JWTClaims
	claims := model.JWTClaims{
		UserID: userID,
		RoleID: roleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(), // Token berlaku selama 12 jam
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserByFullName(db *gorm.DB, fullName string) (*model.Users, error) {
	var user model.Users
	err := db.Where("full_name = ?", fullName).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*model.Users, error) {
	var user model.Users
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByPhoneNumber(db *gorm.DB, phoneNumber string) (*model.Users, error) {
	var user model.Users
	err := db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(db *gorm.DB, userID uint) (*model.Users, error) {
	var user model.Users
	err := db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByRoleID(db *gorm.DB, roleID uint) (*model.Users, error) {
	var user model.Users
	err := db.Where("role_id = ?", roleID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers(db *gorm.DB) ([]model.Users, error) {
	var users []model.Users
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}