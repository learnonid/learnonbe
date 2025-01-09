package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"learnonbe/model" // Ensure import path matches your project structure

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateBook creates a new book record in the database
func CreateBook(db *mongo.Database, book *model.Books) error {
	// Check if book already exists
	collection := db.Collection("books")
	count, err := collection.CountDocuments(context.TODO(), bson.M{"book_name": book.BookName})
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("book already exists")
	}

	// Set book ID and created timestamp
	book.BookID = primitive.NewObjectID()
	book.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Validate the year
	if book.Year <= 2018 || book.Year > time.Now().Year() {
		return fmt.Errorf("invalid year")
	}

	// Validate the store link
	if !strings.HasPrefix(book.StoreLink, "https://") {
		return fmt.Errorf("invalid store link")
	}

	// Insert the book into the database
	_, err = collection.InsertOne(context.TODO(), book)
	if err != nil {
		return err
	}

	return nil
}


// GetBookByID retrieves a book from the database by its ID
func GetBookByID(ctx context.Context, db *mongo.Database, bookID primitive.ObjectID) (*model.Books, error) {
	var book model.Books
	err := db.Collection("books").FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("book with ID %s not found", bookID.Hex())
		}
		return nil, err
	}
	return &book, nil
}

// GetAllBooks retrieves all books from the database
func GetAllBooks(ctx context.Context, db *mongo.Database) ([]model.Books, error) {
	var books []model.Books
	cursor, err := db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book model.Books
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("failed to decode book: %v", err)
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return books, nil
}

// UpdateBook updates a book's information based on its ID
func UpdateBook(ctx context.Context, db *mongo.Database, bookID primitive.ObjectID, updateData bson.M) error {
	collection := db.Collection("books")

	// Update the book in the database
	result, err := collection.UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{"$set": updateData})
	if err != nil {
		return fmt.Errorf("failed to update book: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("book with ID %s not found", bookID.Hex())
	}

	return nil
}

// DeleteBook deletes a book from the database by its ID
func DeleteBook(ctx context.Context, db *mongo.Database, bookID primitive.ObjectID) error {
	collection := db.Collection("books")

	// Delete the book from the database
	result, err := collection.DeleteOne(ctx, bson.M{"_id": bookID})
	if err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("book with ID %s not found", bookID.Hex())
	}

	return nil
}

// GetBooksByAuthor retrieves books by a specific author
func GetBooksByAuthor(ctx context.Context, db *mongo.Database, author string) ([]model.Books, error) {
	var books []model.Books
	filter := bson.M{"author": bson.M{"$eq": author}}
	cursor, err := db.Collection("books").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch books by author: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book model.Books
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("failed to decode book: %v", err)
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return books, nil
}
