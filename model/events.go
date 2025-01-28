package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Struktur untuk acara pelatihan (Events)
type Events struct {
	EventID     primitive.ObjectID `bson:"_id,omitempty" json:"event_id"`            // ID unik di MongoDB
	EventName   string             `bson:"event_name" json:"event_name"`             // Nama acara
	EventType   string             `bson:"event_type" json:"event_type"`             // 'online' atau 'offline'
	EventDate   string             `bson:"event_date" json:"event_date"`             // Tanggal acara
	EventImage  string             `bson:"event_image,omitempty" json:"event_image"` // URL gambar acara
	Location    string             `bson:"location,omitempty" json:"location"`       // Lokasi acara (opsional)
	Price       float64            `bson:"price" json:"price"`                       // Harga reguler
	VIPPrice    float64            `bson:"vip_price" json:"vip_price"`               // Harga VIP
	Description string             `bson:"description,omitempty" json:"description"` // Deskripsi acara (opsional)
	CreatedAt   primitive.DateTime `bson:"created_at,omitempty" json:"created_at"`   // Timestamp saat dokumen dibuat
}
