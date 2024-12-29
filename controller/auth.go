package controller
// LSla7VHMOwgm5STP pw supa
import (
	"learnonbe/model"
	repo "learnonbe/repository"
	
	// "strings"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterAkun(c *fiber.Ctx) error {
	var user model.Users

	// Dapatkan koneksi database dari context
	db := c.Locals("db").(*gorm.DB)

	// Parse body request ke struct User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validasi apakah email sudah terdaftar
	existingUser := new(model.Users)
	if err := db.Where("email = ?", user.Email).First(existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}
	user.Password = string(hashedPassword)

	// Set role ID default (contoh: 2 untuk user biasa)
	if user.RoleID == 0 {
		user.RoleID = 2
	}

	// Simpan user ke database
	if err := repo.CreateUser(db, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Dapatkan koneksi database dari context
	db := c.Locals("db").(*gorm.DB)

	// Parse body request ke struct
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Cari user berdasarkan email
	var user model.Users
	if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Bandingkan password yang diinput dengan yang tersimpan
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := repo.GenerateToken(user.UserID, user.RoleID)
	if err != nil {
		fmt.Println("Error generating token:", err) // Log error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func LoginAdmin(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Dapatkan koneksi database dari context
	db := c.Locals("db").(*gorm.DB)

	// Parse body request ke struct
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Cari user berdasarkan email
	var user model.Users
	if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Periksa apakah user memiliki RoleID untuk admin
	if user.RoleID != 1 { // Misalnya RoleID=1 adalah admin
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied. Admin privileges required.",
		})
	}

	// Bandingkan password yang diinput dengan yang tersimpan
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := repo.GenerateToken(user.UserID, user.RoleID)
	if err != nil {
		fmt.Println("Error generating token:", err) // Log error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Admin login successful",
		"token":   token,
	})
}


func GetProfile(c *fiber.Ctx) error {
    // Ambil klaim dari context
    claims, ok := c.Locals("claims").(*model.JWTClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Unauthorized access",
        })
    }

    // Debug klaim
    fmt.Printf("Claims in GetProfile: UserID=%d, RoleID=%d\n", claims.UserID, claims.RoleID)

    // Ambil UserID dari klaim
    userID := claims.UserID

    // Koneksi ke database
    db := c.Locals("db").(*gorm.DB)

    // Cari user berdasarkan UserID
    user, err := repo.GetUserByID(db, userID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "message": "User not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to retrieve user profile",
            "error":   err.Error(),
        })
    }

    // Hapus informasi sensitif seperti password sebelum mengirim respons
    user.Password = ""

	fmt.Printf("Claims in GetProfile: UserID=%d, RoleID=%d\n", claims.UserID, claims.RoleID)

    // Kembalikan profil pengguna
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "User profile retrieved successfully",
        "user":    user,
    })
}
