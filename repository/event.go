package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"learnonbe/model" // Adjust the import path to your project structure

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateEvent creates a new event in the database
func CreateEvent(db *mongo.Database, event *model.Events) error {
    // Set event ID and created at timestamp
    event.EventID = primitive.NewObjectID()
	event.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    // Insert the event into the database
    collection := db.Collection("events")
    _, err := collection.InsertOne(context.TODO(), event)
    if err != nil {
        return err
    }

    return nil
}

// UploadEventImage handles the upload of an event image and returns the file URL or an error
func UploadEventImage(file *multipart.FileHeader, uploadDir string) (string, error) {
    // Create the upload directory if it doesn't exist
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        fmt.Printf("Error creating upload dir: %v\n", err)
        return "", fmt.Errorf("failed to create upload directory: %v", err)
    }

    // Save the uploaded file to the specified path
    filePath := filepath.Join(uploadDir, file.Filename)
    if err := saveMultipartFile(file, filePath); err != nil {
        fmt.Printf("Error saving file: %v\n", err)
        return "", fmt.Errorf("failed to save file: %v", err)
    }

    // Generate the file URL
    fileURL := fmt.Sprintf("http://localhost:3000/uploads/events/%s", file.Filename)
    fmt.Printf("File URL: %s\n", fileURL)
    return fileURL, nil
}

// saveMultipartFile saves the uploaded file to the given path
func saveMultipartFile(file *multipart.FileHeader, dst string) error {
    fmt.Printf("Saving file: %s to %s\n", file.Filename, dst)
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}

// GetEventByID retrieves an event by its ID
func GetEventByID(ctx context.Context, db *mongo.Database, eventID primitive.ObjectID) (*model.Events, error) {
    var event model.Events
    err := db.Collection("events").FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("event with ID %s not found", eventID.Hex())
        }
        return nil, err
    }
    return &event, nil
}

func GetAllEvents(ctx context.Context, db *mongo.Database) ([]model.Events, error) {
    var events []model.Events
    cursor, err := db.Collection("events").Find(ctx, bson.M{})
    if err != nil {
        return nil, fmt.Errorf("failed to fetch events: %v", err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var event model.Events
        if err := cursor.Decode(&event); err != nil {
            return nil, fmt.Errorf("failed to decode event: %v", err)
        }
        events = append(events, event)
    }

    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %v", err)
    }

    return events, nil
}

func UpdateEvent(ctx context.Context, db *mongo.Database, eventID primitive.ObjectID, updateData bson.M) error {
    collection := db.Collection("events")

    // Update the event in the database
    result, err := collection.UpdateOne(ctx, bson.M{"_id": eventID}, bson.M{"$set": updateData})
    if err != nil {
        return fmt.Errorf("failed to update event: %v", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("event with ID %s not found", eventID.Hex())
    }

    return nil
}

func DeleteEvent(ctx context.Context, db *mongo.Database, eventID primitive.ObjectID) error {
    collection := db.Collection("events")

    // Delete the event from the database
    result, err := collection.DeleteOne(ctx, bson.M{"_id": eventID})
    if err != nil {
        return fmt.Errorf("failed to delete event: %v", err)
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("event with ID %s not found", eventID.Hex())
    }

    return nil
}

func GetEventByDate(ctx context.Context, db *mongo.Database, date time.Time) ([]model.Events, error) {
    var events []model.Events
    filter := bson.M{
        "event_date": bson.M{
            "$eq": date,
        },
    }
    cursor, err := db.Collection("events").Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch events by date: %v", err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var event model.Events
        if err := cursor.Decode(&event); err != nil {
            return nil, fmt.Errorf("failed to decode event: %v", err)
        }
        events = append(events, event)
    }

    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %v", err)
    }

    return events, nil
}

func GetEventByType(ctx context.Context, db *mongo.Database, eventType string) ([]model.Events, error) {
    var events []model.Events
    filter := bson.M{
        "event_type": eventType,
    }
    cursor, err := db.Collection("events").Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch events by type: %v", err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var event model.Events
        if err := cursor.Decode(&event); err != nil {
            return nil, fmt.Errorf("failed to decode event: %v", err)
        }
        events = append(events, event)
    }

    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %v", err)
    }

    return events, nil
}

func GetEventByPrice(ctx context.Context, db *mongo.Database, price float64) ([]model.Events, error) {
    var events []model.Events
    filter := bson.M{
        "price": price,
    }
    cursor, err := db.Collection("events").Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch events by price: %v", err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var event model.Events
        if err := cursor.Decode(&event); err != nil {
            return nil, fmt.Errorf("failed to decode event: %v", err)
        }
        events = append(events, event)
    }

    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %v", err)
    }

    return events, nil
}