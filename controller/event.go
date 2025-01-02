package controller

import (
	"fmt"
	"learnonbe/config"
	"learnonbe/model"
	"learnonbe/repository"
	// "go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/fiber/v2"
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
