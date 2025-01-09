package controller

import (
	"learnonbe/config"
	"learnonbe/model"
	"learnonbe/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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

func UpdateUser(c *fiber.Ctx) error {
    // Parse the user ID from the URL parameter
    userID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid user ID",
            "error":   err.Error(),
        })
    }

    // Parse the request body to get the update data
    var updateData model.Users
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }

    // Convert updateData to bson.M
    update := bson.M{
        "name":  updateData.FullName,
        "email": updateData.Email,
		"phone": updateData.PhoneNumber,
		"password": updateData.Password,
		"status": updateData.Status,
        // Add other fields as needed
    }

    // Get the database connection from context
    db := c.Locals("db").(*mongo.Database)

    // Call the repository function to update the user
    if err := repository.UpdateUser(c.Context(), db, userID, update); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to update user",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "User updated successfully",
    })
}

func DeleteUser(c *fiber.Ctx) error {
    // Parse the user ID from the URL parameter
    userID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid user ID",
            "error":   err.Error(),
        })
    }

    // Get the database connection from context
    db := c.Locals("db").(*mongo.Database)

    // Call the repository function to delete the user
    if err := repository.DeleteUser(c.Context(), db, userID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to delete user",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "User deleted successfully",
    })
}

// Login mengubah koneksi ke MongoDB menggunakan config.MongoClient
func Login(c *fiber.Ctx) error {
    var loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Ambil koneksi database dari config
    db := config.MongoClient.Database("learnon")

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
    db := config.MongoClient.Database("learnon")

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
    db := config.MongoClient.Database("learnon")

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

func LogOut(c *fiber.Ctx) error {
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Logout successful",
    })
}

func GetAllUsers(c *fiber.Ctx) error {
    // Get the database connection from context
    db := config.MongoClient.Database("learnon")

    // Call the repository function to get all users
    users, err := repository.GetAllUsers(c.Context(), db)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to get users",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "users": users,
    })
}   

func GetUserByID(c *fiber.Ctx) error {
    // Parse the user ID from the URL parameter
    userID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid user ID",
            "error":   err.Error(),
        })
    }

    // Get the database connection from context
    db := config.MongoClient.Database("learnon")

    // Call the repository function to get the user by ID
    user, err := repository.GetUserByID(c.Context(), db, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to get user",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "user": user,
    })
}

func GetUserByEmail(c *fiber.Ctx) error {
    // Parse the user email from the URL parameter
    userEmail := c.Params("email")

    // Get the database connection from context
    db := config.MongoClient.Database("learnon")

    // Call the repository function to get the user by email
    user, err := repository.GetUserByEmail(c.Context(), db, userEmail)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to get user",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "user": user,
    })
}