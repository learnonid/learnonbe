package controller

import (
	"learnonbe/config" // Import package config
	"learnonbe/model"
	"learnonbe/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAkun mengubah koneksi ke MongoDB menggunakan config.MongoClient
func RegisterAkun(c *fiber.Ctx) error {
    var user model.Users

    // Ambil koneksi database dari config
    db := config.MongoClient.Database("learnon")

    // Parse body request ke struct User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // Validasi apakah email sudah terdaftar
    collection := db.Collection("users")
    count, err := collection.CountDocuments(c.Context(), bson.M{"email": user.Email})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to check email existence",
            "error":   err.Error(),
        })
    }
    if count > 0 {
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

    // Set role ID default ke 2 (customer) jika tidak ada input RoleID
    if user.RoleID == 0 {  // Jika RoleID kosong, berarti pengguna tidak memberikan role
        user.RoleID = 2  // Set default ke customer
    }

    // Simpan user ke database
    user.UserID = primitive.NewObjectID()
    _, err = collection.InsertOne(c.Context(), user)
    if err != nil {
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


// Login mengubah koneksi ke MongoDB menggunakan config.MongoClient
func Login(c *fiber.Ctx) error {
    var loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Ambil koneksi database dari config
    db := config.MongoClient.Database("your_database_name")

    // Parse body request ke struct
    if err := c.BodyParser(&loginRequest); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // Cari user berdasarkan email
    collection := db.Collection("users")
    var user model.Users
    err := collection.FindOne(c.Context(), bson.M{"email": loginRequest.Email}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid email or password",
        })
    }

    // Verifikasi password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid email or password",
        })
    }

    // Buat token JWT (contoh sederhana)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.UserID
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

    t, err := token.SignedString([]byte("secret"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to generate token",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
        "token":   t,
    })
}

// LoginAdmin mengubah koneksi ke MongoDB menggunakan config.MongoClient
func LoginAdmin(c *fiber.Ctx) error {
    var loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Ambil koneksi database dari config
    db := config.MongoClient.Database("your_database_name")

    // Parse body request ke struct
    if err := c.BodyParser(&loginRequest); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // Cari user berdasarkan email
    collection := db.Collection("users")
    var user model.Users
    err := collection.FindOne(c.Context(), bson.M{"email": loginRequest.Email, "role": "admin"}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid email or password",
        })
    }

    // Verifikasi password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid email or password",
        })
    }

    // Buat token JWT (contoh sederhana)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.UserID
    claims["role"] = "admin"
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

    t, err := token.SignedString([]byte("secret"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to generate token",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
        "token":   t,
    })
}

// GetProfile mengubah koneksi ke MongoDB menggunakan config.MongoClient
func GetProfile(c *fiber.Ctx) error {
    claims, ok := c.Locals("claims").(*model.JWTClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Unauthorized access",
        })
    }

    // Ambil UserID dari klaim
    userID := claims.UserID

    // Koneksi ke database
    db := config.MongoClient.Database("your_database_name")

    // Gunakan repository untuk mencari user berdasarkan UserID
    user, err := repository.GetUserByID(c.Context(), db, userID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User not found",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "user": user,
    })
}
