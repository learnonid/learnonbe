package controller

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/learnonid/learnonbe/config"
	"github.com/learnonid/learnonbe/model"
	"github.com/learnonid/learnonbe/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostBookPayment(c *fiber.Ctx) error {
	// Ambil data dari form-data
	userID := c.FormValue("user_id")
	userName := c.FormValue("full_name")
	bookID := c.FormValue("book_id")
	bookName := c.FormValue("book_name") // Nama buku
	priceStr := c.FormValue("price")     // Harga
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid price format",
		})
	}

	// Validasi konversi userID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Ambil file pembayaran
	paymentFile, err := c.FormFile("payment_receipt")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment receipt file is required",
		})
	}

	// Baca konten file pembayaran
	paymentFileContent, err := paymentFile.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to open payment receipt file",
		})
	}
	defer paymentFileContent.Close()

	// Konversi file pembayaran ke base64
	paymentContent, err := ioutil.ReadAll(paymentFileContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read payment receipt file content",
		})
	}
	paymentBase64Content := base64.StdEncoding.EncodeToString(paymentContent)

	// Buat nama file untuk pembayaran
	paymentFileName := fmt.Sprintf("receipts/%s_%s_%d_%s", userID, bookID, time.Now().Unix(), paymentFile.Filename)

	// Upload file pembayaran ke GitHub
	err = repository.UploadToGithub(paymentFileName, paymentBase64Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Simpan data ke MongoDB dalam koleksi "book_payment"
	bookPayment := model.BookPayment{
		UserID:           objectID, // userID yang sudah di-convert ke ObjectID
		UserName:         userName,
		BookID:           primitive.NewObjectID(),
		BookName:         bookName,
		Price:            price,
		PaymentReceipt:   fmt.Sprintf("https://github.com/learnonid/uploads/blob/main/%s", paymentFileName),
		PaymentDate: primitive.NewDateTimeFromTime(time.Now()),
	}

	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "book_payment"
	collection := client.Database("learnon").Collection("book_payment")

	// Simpan data
	_, err = collection.InsertOne(c.Context(), bookPayment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save payment data",
		})
	}

	// Respons berhasil
	return c.JSON(fiber.Map{
		"message": "File uploaded and payment saved successfully",
		"url":     fmt.Sprintf("https://github.com/learnonid/uploads/blob/main/%s", paymentFileName),
		"data": fiber.Map{
			"user_id":   userID,
			"full_name": userName,
			"book_id":   bookID,
			"book_name": bookName,
			"price":     price,
		},
	})
}

func GetAllBookPayment(c *fiber.Ctx) error {
	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "book_payment"
	collection := client.Database("learnon").Collection("book_payment")

	// Ambil semua data payment
	cursor, err := collection.Find(c.Context(), primitive.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch registration data",
		})
	}
	defer cursor.Close(c.Context())

	// Loop semua data payment
	var bookPayments []model.BookPayment
	for cursor.Next(c.Context()) {
		var payment model.BookPayment
		err := cursor.Decode(&payment)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode payment data",
			})
		}
		bookPayments = append(bookPayments, payment)
	}

	// Respons data payment
	return c.JSON(bookPayments)
}

func GetBookPaymentByUserID(c *fiber.Ctx) error {
	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "book_payment"
	collection := client.Database("learnon").Collection("book_payment")

	// Ambil user ID dari parameter URL
	userIDStr := c.Params("id")

	// Konversi string userID menjadi ObjectId
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Ambil data payment berdasarkan user ID
	cursor, err := collection.Find(c.Context(), primitive.D{{Key: "user_id", Value: userID}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch payment data",
		})
	}
	defer cursor.Close(c.Context())

	// Loop semua data payment
	var bookPayments []model.BookPayment
	for cursor.Next(c.Context()) {
		var payment model.BookPayment
		err := cursor.Decode(&payment)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode payment data",
			})
		}
		bookPayments = append(bookPayments, payment)
	}

	// Respons data registrasi
	if len(bookPayments) == 0 {
		return c.JSON(fiber.Map{
			"message": "No payment found for this user",
		})
	}
	return c.JSON(bookPayments)
}

func GetBookPaymentByID(c *fiber.Ctx) error {
	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "book_payment"
	collection := client.Database("learnon").Collection("book_payment")

	// Ambil ID dari parameter URL
	objectIDStr := c.Params("id")

	// Konversi string ID menjadi ObjectId
	objectID, err := primitive.ObjectIDFromHex(objectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	// Ambil data registrasi berdasarkan ID
	filter := bson.M{"_id": objectID}
	var payment model.BookPayment
	err = collection.FindOne(c.Context(), filter).Decode(&payment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch payment data",
		})
	}

	// Respons data registrasi
	return c.JSON(payment)
}