package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Books struct {
	BookID     primitive.ObjectID `bson:"_id,omitempty" json:"book_id"` // ID unik di MongoDB
	BookName   string             `bson:"book_name" json:"book_name"`   // Nama buku
	Author     string             `bson:"author" json:"author"`         // Penulis buku
	Publisher  string             `bson:"publisher" json:"publisher"`   // Penerbit buku
	Year       int                `bson:"year" json:"year"`             // Tahun terbit
	ISBN       string             `bson:"isbn" json:"isbn"`             // Nomor ISBN
	Price      float64            `bson:"price" json:"price"`           // Harga buku
	StoreLink  string             `bson:"store_link,omitempty" json:"store_link"` // URL toko buku
	CreatedAt  primitive.DateTime `bson:"created_at,omitempty" json:"created_at"` // Timestamp saat dokumen dibuat
}