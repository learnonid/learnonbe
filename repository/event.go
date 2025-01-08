package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"learnonbe/model" // Pastikan path impor sesuai dengan struktur proyek Anda

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateEvents membuat acara kursus baru dalam database
func CreateEvent(db *mongo.Database, event *model.Events) error {
	// Set ID acara dan timestamp pembuatan
	event.EventID = primitive.NewObjectID()
	event.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Masukkan acara ke dalam database
	collection := db.Collection("events")
	_, err := collection.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}

	return nil
}

// UploadEventImage menangani unggahan gambar acara dan mengembalikan URL file atau error
func UploadEventImage(file *multipart.FileHeader, uploadDir string) (string, error) {
	// Buat direktori upload jika belum ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating upload dir: %v\n", err)
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Simpan file yang diunggah ke path yang ditentukan
	filePath := filepath.Join(uploadDir, file.Filename)
	if err := saveMultipartFile(file, filePath); err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Hasilkan URL file
	fileURL := fmt.Sprintf("http://localhost:3000/uploads/courses/%s", file.Filename)
	fmt.Printf("File URL: %s\n", fileURL)
	return fileURL, nil
}

// saveMultipartFile menyimpan file yang diunggah ke path yang diberikan
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

// GetEventsByID mengambil acara kursus berdasarkan ID-nya
func GetEventsByID(ctx context.Context, db *mongo.Database, eventID primitive.ObjectID) (*model.Events, error) {
	var event model.Events
	err := db.Collection("events").FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("course event with ID %s not found", eventID.Hex())
		}
		return nil, err
	}
	return &event, nil
}

// GetAllEventss mengambil semua acara kursus yang ada
func GetAllEvents(ctx context.Context, db *mongo.Database) ([]model.Events, error) {
	var events []model.Events
	cursor, err := db.Collection("events").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch course events: %v", err)
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

// UpdateEvents memperbarui data acara kursus berdasarkan ID
func UpdateEvents(ctx context.Context, db *mongo.Database, eventID primitive.ObjectID, updateData bson.M) error {
	collection := db.Collection("events")

	// Perbarui acara kursus dalam database
	result, err := collection.UpdateOne(ctx, bson.M{"_id": eventID}, bson.M{"$set": updateData})
	if err != nil {
		return fmt.Errorf("failed to update event: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("course event with ID %s not found", eventID.Hex())
	}

	return nil
}

// DeleteEvents menghapus acara kursus berdasarkan ID
func DeleteEvents(ctx context.Context, db *mongo.Database, eventID primitive.ObjectID) error {
	collection := db.Collection("events")

	// Hapus acara kursus dari database
	result, err := collection.DeleteOne(ctx, bson.M{"_id": eventID})
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("course event with ID %s not found", eventID.Hex())
	}

	return nil
}

// GetEventsByDate mengambil acara kursus berdasarkan tanggal
func GetEventsByDate(ctx context.Context, db *mongo.Database, date time.Time) ([]model.Events, error) {
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

// GetEventsByType mengambil acara kursus berdasarkan jenisnya
func GetEventsByType(ctx context.Context, db *mongo.Database, eventType string) ([]model.Events, error) {
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

// GetEventsByType Online mengambil acara kursus online
func GetEventsByTypeOnline(ctx context.Context, db *mongo.Database) ([]model.Events, error) {
	return GetEventsByType(ctx, db, "online")
}

// GetEventsByTypeOffline mengambil acara kursus offline
func GetEventsByTypeOffline(ctx context.Context, db *mongo.Database) ([]model.Events, error) {
	return GetEventsByType(ctx, db, "offline")
}

// GetEventsByPrice mengambil acara kursus berdasarkan harga
func GetEventsByPrice(ctx context.Context, db *mongo.Database, price float64) ([]model.Events, error) {
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
