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

func getStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func PostUploadPayment(c *fiber.Ctx) error {
	// Ambil data dari form-data
	userID := c.FormValue("user_id")
	userName := c.FormValue("full_name")
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
		UserName:         userName,
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
			"full_name":       userName,
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

func GetUERegistrationByID(c *fiber.Ctx) error {
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
	var registration model.UserEventRegistration
	err = collection.FindOne(c.Context(), filter).Decode(&registration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch registration data",
		})
	}

	// Respons data registrasi
	return c.JSON(registration)
}

func UpdateEventRegistration(c *fiber.Ctx) error {
	// Ambil _id dari parameter URL
	objectIDStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(objectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format. Make sure the ID is a valid ObjectID.",
		})
	}

	// Ambil data dari form-data
	userID := c.FormValue("user_id")
	eventID := c.FormValue("event_id")
	eventName := c.FormValue("event_name")
	status := c.FormValue("status")
	priceStr := c.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid price format",
		})
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

	// Filter berdasarkan _id
	filter := bson.M{
		"_id": objectID,
	}

	// Cari dokumen di database untuk memeriksa kondisi file
	var existingData bson.M
	if err := collection.FindOne(c.Context(), filter).Decode(&existingData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch existing data",
		})
	}

	// Siapkan variabel untuk file sertifikat dan materi
	var sertifikatURL string
	var materiURL string

	// Cek apakah sertifikat_file adalah link atau form-data
	if existingData["sertifikat_file"] == nil || existingData["sertifikat_file"] == "" {
		// Ambil file sertifikat dari form-data
		sertifikatFile, err := c.FormFile("sertifikat_file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Sertifikat file is required",
			})
		}

		// Baca konten file sertifikat
		sertifikatFileContent, err := sertifikatFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to open sertifikat file",
			})
		}
		defer sertifikatFileContent.Close()

		// Konversi file sertifikat ke base64
		sertifikatContent, err := ioutil.ReadAll(sertifikatFileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read sertifikat file content",
			})
		}
		sertifikatBase64 := base64.StdEncoding.EncodeToString(sertifikatContent)

		// Buat nama file untuk sertifikat
		sertifikatFileName := fmt.Sprintf("sertifikat/%s_%s_%d_%s", userID, eventID, time.Now().Unix(), sertifikatFile.Filename)

		// Upload file sertifikat ke GitHub
		err = repository.UploadToGithub(sertifikatFileName, sertifikatBase64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		sertifikatURL = fmt.Sprintf("https://github.com/learnonid/uploads/blob/main/%s", sertifikatFileName)
	} else {
		// Gunakan link yang ada di database
		sertifikatURL = existingData["sertifikat_file"].(string)
	}

	// Cek apakah materi_file adalah link atau form-data
	if existingData["materi_file"] == nil || existingData["materi_file"] == "" {
		// Ambil file materi dari form-data
		materiFile, err := c.FormFile("materi_file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Materi file is required",
			})
		}

		// Baca konten file materi
		materiFileContent, err := materiFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to open materi file",
			})
		}
		defer materiFileContent.Close()

		// Konversi file materi ke base64
		materiContent, err := ioutil.ReadAll(materiFileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read materi file content",
			})
		}
		materiBase64 := base64.StdEncoding.EncodeToString(materiContent)

		// Buat nama file untuk materi
		materiFileName := fmt.Sprintf("materi/%s_%s_%d_%s", userID, eventID, time.Now().Unix(), materiFile.Filename)

		// Upload file materi ke GitHub
		err = repository.UploadMateri(materiFileName, materiBase64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		materiURL = fmt.Sprintf("https://github.com/learnonid/uploads/blob/main/%s", materiFileName)
	} else {
		// Gunakan link yang ada di database
		materiURL = existingData["materi_file"].(string)
	}

	// Update data di MongoDB
	update := bson.M{
		"$set": bson.M{
			"event_name":      eventName,
			"status":          status,
			"price":           price,
			"sertifikat_file": sertifikatURL,
			"materi_file":     materiURL,
		},
	}

	_, err = collection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update registration data",
		})
	}

	// Respons berhasil
	return c.JSON(fiber.Map{
		"message": "Registration data updated successfully",
		"data": fiber.Map{
			"user_id":         userID,
			"event_id":        eventID,
			"event_name":      eventName,
			"status":          status,
			"price":           price,
			"sertifikat_file": sertifikatURL,
			"materi_file":     materiURL,
		},
	})
}
