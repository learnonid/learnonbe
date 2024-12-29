package controller

import (
	"learnonbe/model"
	"learnonbe/repository"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateEvent(c *fiber.Ctx) error {
	var event model.Events

	// Parse JSON body into the Event struct
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

	// Get the database connection from context
	db := c.Locals("db").(*gorm.DB)

	// Call the repository function to create the event
	if err := repository.CreateEvent(db, &event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create event",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"event":   event,
	})
}

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
        "message": "Image uploaded successfully",
        "file_url": fileURL,
    })
}
