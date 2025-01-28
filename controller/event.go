package controller

import (
	"fmt"
	"log"
	"strconv"

	"github.com/learnonid/learnonbe/config"
	"github.com/learnonid/learnonbe/model"
	"github.com/learnonid/learnonbe/repository"

	// "go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEvent handles creating an event in the database
func CreateEvent(c *fiber.Ctx) error {
	var event model.Events

	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Parse body request to struct Event
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Handle file upload
	file, err := c.FormFile("event_image")
	if err == nil { // File upload is optional
		fileURL, err := repository.UploadEventImage(file, "./uploads/events")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to upload event image",
				"error":   err.Error(),
			})
		}
		event.EventImage = fileURL
	}

	// Call the repository function to create the event in MongoDB
	if err := repository.CreateEvent(db, &event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create event",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"event":   event,
	})
}

// UploadEventImageHandler handles event image upload
func UploadEventImageHandler(c *fiber.Ctx) error {
	fmt.Printf("Headers: %v\n", c.GetReqHeaders()) // Debug header request
	file, err := c.FormFile("event_image")
	if err != nil {
		fmt.Printf("FormFile error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to retrieve file",
			"error":   err.Error(),
		})
	}

	fileURL, err := repository.UploadEventImage(file, "./uploads/events")
	if err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
	}

	fmt.Printf("File uploaded to: %s\n", fileURL)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Image uploaded successfully",
		"file_url": fileURL,
	})
}

// GetEvents retrieves all events from the database
func GetEvents(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Call the repository function to get all events
	events, err := repository.GetAllEvents(c.Context(), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch events",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Events retrieved successfully",
		"events":  events,
	})
}

// GetEventByID retrieves an event by ID from the database
func GetEventByID(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Get the event ID from the URL parameter
	eventID := c.Params("id")
	// Convert string eventID to ObjectID
	eventIDObj, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid event ID format",
			"error":   err.Error(),
		})
	}

	// Call the repository function to get the event by ID
	event, err := repository.GetEventsByID(c.Context(), db, eventIDObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch event",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event retrieved successfully",
		"event":   event,
	})
}

// GetEventsByType retrieves events by their type (online or offline)
func GetEventsByType(c *fiber.Ctx) error {
	// Ambil event_type dari query parameter
	eventType := c.Query("type")
	if eventType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "event_type is required",
		})
	}

	// Akses database dari konfigurasi
	db := config.MongoClient.Database("learnon")

	// Panggil fungsi repository
	events, err := repository.GetEventsByType(c.Context(), db, eventType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch events",
			"error":   err.Error(),
		})
	}

	if len(events) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No events found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Events retrieved successfully",
		"events":  events,
	})
}

func GetEventByTypeOnline(c *fiber.Ctx) error {
	// Tetapkan eventType sebagai "online"
	eventType := "online"

	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Panggil repository untuk mendapatkan events berdasarkan type
	events, err := repository.GetEventsByType(c.Context(), db, eventType)

	if err != nil {
		// Jika terjadi error saat mengambil data
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch online events",
			"error":   err.Error(),
		})
	}

	// Jika tidak ada event ditemukan
	if len(events) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No online events found",
		})
	}

	// Jika berhasil, kirimkan response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Online events retrieved successfully",
		"events":  events,
	})
}

// GetEventByTypeOffline retrieves all offline events from the database
func GetEventByTypeOffline(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Call the repository function to get all offline events
	events, err := repository.GetEventsByTypeOffline(c.Context(), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch offline events",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Offline events retrieved successfully",
		"events":  events,
	})
}

// EditEvent handles updating an event in the database
func EditEvent(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Get the event ID from the URL parameter
	eventID := c.Params("id")
	eventIDObj, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid event ID format",
			"error":   err.Error(),
		})
	}

	// Parse form-data manually
	eventName 	:= c.FormValue("event_name")
	eventType 	:= c.FormValue("event_type")
	eventDate 	:= c.FormValue("event_date")
	location 	:= c.FormValue("location")
	description := c.FormValue("description")

	// Parse price and vipPrice as float64
	priceStr 	:= c.FormValue("price")
	vipPriceStr := c.FormValue("vip_price")

	// Handle file upload
	var fileURL string
	file, err := c.FormFile("event_image")
	if err == nil { // File upload is optional
		fileURL, err = repository.UploadEventImage(file, "./uploads/events")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to upload event image",
				"error":   err.Error(),
			})
		}
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println("Failed to parse price:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid price format",
			"error":   err.Error(),
		})
	}

	vipPrice, err := strconv.ParseFloat(vipPriceStr, 64)
	if err != nil {
		log.Println("Failed to parse vip_price:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid vip_price format",
			"error":   err.Error(),
		})
	}

	// Prepare the update data
	updateData := bson.M{
		"event_name":  eventName,
		"event_type":  eventType,
		"event_date":  eventDate,
		"location":    location,
		"description": description,
		"price":       price,
		"vip_price":   vipPrice,
	}
	if fileURL != "" {
		updateData["event_image"] = fileURL
	}

	// Call the repository function to update the event in MongoDB
	if err := repository.UpdateEvents(c.Context(), db, eventIDObj, updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update event",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event updated successfully",
		"event":   updateData,
	})
}

// DeleteEvent handles deleting an event from the database
func DeleteEvent(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Get the event ID from the URL parameter
	eventID := c.Params("id")
	// Convert string eventID to ObjectID
	eventIDObj, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid event ID format",
			"error":   err.Error(),
		})
	}

	// Call the repository function to delete the event from MongoDB
	if err := repository.DeleteEvents(c.Context(), db, eventIDObj); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete event",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}
