package controller

import (
	"context"
	// "fmt"
	"learnonbe/model"
	"learnonbe/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create Roles
func CreateRole(c *fiber.Ctx) error {
	var role model.Roles

	// Parse body request to struct Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Get the database connection from context
	db := c.Locals("db").(*mongo.Database)

	// Call the repository function to create the role in MongoDB
	if err := repository.CreateRole(context.TODO(), db, &role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create role",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Role created successfully",
		"role":    role,
	})
}