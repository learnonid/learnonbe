package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struktur untuk registrasi pengguna pada acara (User_Event_Registrations)
type UserEventRegistration struct {
	RegistrationID   primitive.ObjectID `bson:"_id,omitempty" json:"registration_id"`  // ID unik registrasi
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`               // ID pengguna (referensi ke Users)
	EventID          primitive.ObjectID `bson:"event_id" json:"event_id"`             // ID acara (referensi ke Events)
	Status           string             `bson:"status" json:"status"`                 // 'biasa', 'vip', 'reject'
	PaymentReceipt   string             `bson:"payment_receipt,omitempty" json:"payment_receipt"` // Lokasi bukti pembayaran
	MateriFile       string             `bson:"materi_file,omitempty" json:"materi_file"`         // Lokasi file materi untuk user VIP
	SertifikatFile   string             `bson:"sertifikat_file,omitempty" json:"sertifikat_file"` // Lokasi file sertifikat untuk user VIP
	RegistrationDate primitive.DateTime `bson:"registration_date,omitempty" json:"registration_date"` // Tanggal registrasi
}
