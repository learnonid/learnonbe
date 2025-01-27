package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struktur untuk registrasi pengguna pada acara (User_Event_Registrations)
type BookPayment struct {
	RegistrationID primitive.ObjectID `bson:"_id,omitempty" json:"registration_id"`             // ID unik registrasi
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`                           // ID pengguna (referensi ke Users)
	BookID         primitive.ObjectID `bson:"book_id" json:"book_id"`                           // ID acara (referensi ke Events)
	UserName       string             `bson:"full_name" json:"full_name"`                       // Nama user
	BookName       string             `bson:"book_name" json:"book_name"`                       // Nama acara
	Price          float64            `bson:"price" json:"price"`                               // Harga
	PaymentReceipt string             `bson:"payment_receipt,omitempty" json:"payment_receipt"` // Lokasi bukti pembayaran
	PaymentDate    primitive.DateTime `bson:"payment_date,omitempty" json:"payment_date"`       // Tanggal registrasi
}
