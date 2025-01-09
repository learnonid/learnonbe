package controller

import (
	// "fmt"
	"learnonbe/config"
	"learnonbe/model"
	"learnonbe/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateBook handles creating a book in the database
func CreateBook(c *fiber.Ctx) error {
	var book model.Books

	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Parse body request to struct Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}


	// Call the repository function to create the book in MongoDB
	if err := repository.CreateBook(db, &book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create book",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book created successfully",
		"book":    book,
	})
}

// GetBooks retrieves all books from the database
func GetBooks(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Call the repository function to get all books
	books, err := repository.GetAllBooks(c.Context(), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch books",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Books retrieved successfully",
		"books":   books,
	})
}

// GetBookByID retrieves a book by ID from the database
func GetBookByID(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Get the book ID from the URL parameter
	bookID := c.Params("id")
	// Convert string bookID to ObjectID
	bookIDObj, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid book ID format",
			"error":   err.Error(),
		})
	}

	// Call the repository function to get the book by ID
	book, err := repository.GetBookByID(c.Context(), db, bookIDObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch book",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book retrieved successfully",
		"book":    book,
	})
}

// EditBook handles updating a book in the database
func EditBook(c *fiber.Ctx) error {
	var book model.Books

	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Parse body request to struct Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Get the book ID from the URL parameter
	bookID := c.Params("id")
	// Convert string bookID to ObjectID
	bookIDObj, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid book ID format",
			"error":   err.Error(),
		})
	}

	// Prepare the update data (book details to be updated)
	updateData := bson.M{
		"book_name":    book.BookName,
		"author":       book.Author,
		"publisher":    book.Publisher,
		"year":         book.Year,
		"isbn":         book.ISBN,
		"price":        book.Price,
		"store_link":   book.StoreLink,
		"created_at":   book.CreatedAt,
	}

	// Call the repository function to update the book in MongoDB
	if err := repository.UpdateBook(c.Context(), db, bookIDObj, updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update book",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book updated successfully",
		"book":    book,
	})
}

// DeleteBook handles deleting a book from the database
func DeleteBook(c *fiber.Ctx) error {
	// Get the database connection from config
	db := config.MongoClient.Database("learnon")

	// Get the book ID from the URL parameter
	bookID := c.Params("id")
	// Convert string bookID to ObjectID
	bookIDObj, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid book ID format",
			"error":   err.Error(),
		})
	}

	// Call the repository function to delete the book from MongoDB
	if err := repository.DeleteBook(c.Context(), db, bookIDObj); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete book",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book deleted successfully",
	})
}
