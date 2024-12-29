package repository

import (
	"learnonbe/model"
	"learnonbe/utils"

	"os"
	"fmt"
	"io"

	"gorm.io/gorm"
	"path/filepath"
	"mime/multipart"
)

func CreateEvent(db *gorm.DB, event *model.Events) error {
	// Generate EventID secara random
	event.EventID = utils.GenerateRandomID(1, 10000)

	// Simpan event
	err := db.Create(&event).Error
	if err != nil {
		return fmt.Errorf("gagal membuat event: %v", err)
	}

	return nil
}

// UploadEventImage handles the upload of an event image and returns the file URL or an error
func UploadEventImage(file *multipart.FileHeader, uploadDir string) (string, error) {
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        fmt.Printf("Error creating upload dir: %v\n", err)
        return "", fmt.Errorf("failed to create upload directory: %v", err)
    }

    filePath := filepath.Join(uploadDir, file.Filename)
    if err := saveMultipartFile(file, filePath); err != nil {
        fmt.Printf("Error saving file: %v\n", err)
        return "", fmt.Errorf("failed to save file: %v", err)
    }

    fileURL := fmt.Sprintf("http://localhost:3000/uploads/events/%s", file.Filename)
    fmt.Printf("File URL: %s\n", fileURL)
    return fileURL, nil
}


// saveMultipartFile saves the uploaded file to the given path
func saveMultipartFile(file *multipart.FileHeader, dst string) error {
    fmt.Printf("Saving file: %s to %s\n", file.Filename, dst)
    src, err := file.Open()
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return err
    }
    defer src.Close()

    out, err := os.Create(dst)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    if err != nil {
        fmt.Printf("Error copying file: %v\n", err)
    }
    return err
}

func GetEventByID(db *gorm.DB, eventID int) (*model.Events, error) {
	event := new(model.Events)
	err := db.Where("event_id = ?", eventID).First(&event).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil event: %v", err)
	}

	return event, nil
}

func GetAllEvents(db *gorm.DB) ([]model.Events, error) {
	var events []model.Events
	err := db.Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil semua event: %v", err)
	}

	return events, nil
}

func UpdateEvent(db *gorm.DB, event *model.Events) error {
	err := db.Save(&event).Error
	if err != nil {
		return fmt.Errorf("gagal mengupdate event: %v", err)
	}

	return nil
}

func DeleteEvent(db *gorm.DB, eventID int) error {
	err := db.Delete(&model.Events{}, eventID).Error
	if err != nil {
		return fmt.Errorf("gagal menghapus event: %v", err)
	}

	return nil
}

func GetEventByDate(db *gorm.DB, date string) ([]model.Events, error) {
	var events []model.Events
	err := db.Where("event_date = ?", date).Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil event berdasarkan tanggal: %v", err)
	}

	return events, nil
}

func GetEventByType(db *gorm.DB, eventType string) ([]model.Events, error) {
	var events []model.Events
	err := db.Where("event_type = ?", eventType).Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil event berdasarkan tipe: %v", err)
	}

	return events, nil
}

func GetEventByPrice(db *gorm.DB, price float64) ([]model.Events, error) {
	var events []model.Events
	err := db.Where("price = ?", price).Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil event berdasarkan harga: %v", err)
	}

	return events, nil
}