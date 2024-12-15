package controller

import (
	"learnonbe/model"
	repo "learnonbe/repository"
	
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterAkun(c *fiber.Ctx) error {
	var user model.Users

	// Koneksi ke database
	db := c.Locals("db").(*gorm.DB)

	// Parse body request ke struct User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal memproses request",
			"error":   err.Error(),
		})
	}

	// Validasi email
	if err := repo.ValidateEmail(user.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email tidak valid",
			"error":   err.Error(),
		})
	}

	// Validasi phone number
	if err := repo.ValidatePhoneNumber(user.PhoneNumber); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nomor telepon tidak valid",
			"error":   err.Error(),
		})
	}

	// Cek apakah email sudah terdaftar
	if err := db.Where("email = ?", user.Email).First(&model.Users{}).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email sudah terdaftar",
		})
	}

	// Cek apakah nomor telepon sudah terdaftar
	if err := db.Where("phone_number = ?", user.PhoneNumber).First(&model.Users{}).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nomor telepon sudah terdaftar",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghash password",
			"error":   err.Error(),
		})
	}
	user.Password = string(hashedPassword)

	// Menyimpan data user ke database
	if err := repo.CreateUser(db, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Registrasi gagal",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registrasi berhasil, Silahkan login",
	})
}