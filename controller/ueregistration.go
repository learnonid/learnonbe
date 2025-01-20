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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func PostUploadPayment(c *fiber.Ctx) error {
	// Ambil data dari form-data
	userID := c.FormValue("user_id")
	eventID := c.FormValue("event_id")
	eventName := c.FormValue("event_name") // Nama acara
	status := c.FormValue("status")        // Status (regular, VIP, etc.)
	priceStr := c.FormValue("price")       // Harga
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
	paymentFileName := fmt.Sprintf("receipts/%s_%s_%d_%s", userID, eventID, time.Now().Unix(), paymentFile.Filename)

	// Upload file pembayaran ke GitHub
	err = repository.UploadToGithub(paymentFileName, paymentBase64Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil file materi jika ada
	materiFile := c.FormValue("materi_file") // Ini bisa berupa URL atau nama file materi
	var materiFileValue *string
	if materiFile != "" {
		materiFileValue = &materiFile
	}

	// Ambil file sertifikat jika ada
	sertifikatFile := c.FormValue("sertifikat_file") // Ini bisa berupa URL atau nama file sertifikat
	var sertifikatFileValue *string
	if sertifikatFile != "" {
		sertifikatFileValue = &sertifikatFile
	}

	// Simpan data ke MongoDB dalam koleksi "ueregist"
	registrasi := model.UserEventRegistration{
		UserID:           objectID, // userID yang sudah di-convert ke ObjectID
		EventID:          primitive.NewObjectID(),
		EventName:        eventName,
		Status:           status,
		Price:            price,
		PaymentReceipt:   fmt.Sprintf("https://github.com/learnonid/uploads/blob/main/%s", paymentFileName),
		RegistrationDate: primitive.NewDateTimeFromTime(time.Now()),
		MateriFile:       getStringValue(materiFileValue),     // File materi jika ada
		SertifikatFile:   getStringValue(sertifikatFileValue), // File sertifikat jika ada
	}

	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "ueregist"
	collection := client.Database("learnon").Collection("ueregist")

	// Simpan data
	_, err = collection.InsertOne(c.Context(), registrasi)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save registration data",
		})
	}

	// Respons berhasil
	return c.JSON(fiber.Map{
		"message": "File uploaded and registration saved successfully",
		"url":     fmt.Sprintf("https://github.com/learnonid/uploads/blob/main/%s", paymentFileName),
		"data": fiber.Map{
			"user_id":         userID,
			"event_id":        eventID,
			"event_name":      eventName,
			"status":          status,
			"price":           price,
			"materi_file":     materiFileValue,
			"sertifikat_file": sertifikatFileValue,
		},
	})
}

func GetAllUERegistration(c *fiber.Ctx) error {
	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "ueregist"
	collection := client.Database("learnon").Collection("ueregist")

	// Ambil semua data registrasi
	cursor, err := collection.Find(c.Context(), primitive.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch registration data",
		})
	}
	defer cursor.Close(c.Context())

	// Loop semua data registrasi
	var registrations []model.UserEventRegistration
	for cursor.Next(c.Context()) {
		var registration model.UserEventRegistration
		err := cursor.Decode(&registration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode registration data",
			})
		}
		registrations = append(registrations, registration)
	}

	// Respons data registrasi
	return c.JSON(registrations)
}

func GetUERegistrationByUserID(c *fiber.Ctx) error {
	// Ambil koneksi MongoDB
	client := config.GetMongoClient()
	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to the database",
		})
	}
	defer client.Disconnect(c.Context())

	// Pilih koleksi "ueregist"
	collection := client.Database("learnon").Collection("ueregist")

	// Ambil user ID dari parameter URL
	userIDStr := c.Params("id")

	// Konversi string userID menjadi ObjectId
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Ambil data registrasi berdasarkan user ID
	cursor, err := collection.Find(c.Context(), primitive.D{{Key: "user_id", Value: userID}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch registration data",
		})
	}
	defer cursor.Close(c.Context())

	// Loop semua data registrasi
	var registrations []model.UserEventRegistration
	for cursor.Next(c.Context()) {
		var registration model.UserEventRegistration
		err := cursor.Decode(&registration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode registration data",
			})
		}
		registrations = append(registrations, registration)
	}

	// Respons data registrasi
	if len(registrations) == 0 {
		return c.JSON(fiber.Map{
			"message": "No registrations found for this user",
		})
	}
	return c.JSON(registrations)
}
