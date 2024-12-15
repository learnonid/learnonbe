package model

// Struktur untuk acara pelatihan (Events)
type Event struct {
	EventID     uint    `gorm:"primaryKey;autoIncrement" json:"event_id"`
	EventName   string  `gorm:"not null" json:"event_name"`
	EventType   string  `gorm:"not null" json:"event_type"`         // 'online' atau 'offline'
	EventDate   string  `gorm:"not null" json:"event_date"`         // Format Datetime sebagai string
	Location    string  `gorm:"size:255" json:"location,omitempty"` // Hanya untuk acara offline
	Price       float64 `gorm:"not null" json:"price"`
	VIPPrice    float64 `gorm:"not null" json:"vip_price"`
	Description string  `gorm:"type:text" json:"description"` // Deskripsi acara
}
