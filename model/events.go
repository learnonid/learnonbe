package model

// Struktur untuk acara pelatihan (Events)
type Events struct {
	EventID     uint    `gorm:"primaryKey;autoIncrement" json:"event_id"`
	EventName   string  `gorm:"not null" json:"event_name"`
	EventType   string  `gorm:"not null" json:"event_type"` // 'online' atau 'offline'
	EventDate   string  `gorm:"not null" json:"event_date"`
	EventImage  string  `gorm:"size:255" json:"event_image"`
	Location    string  `gorm:"size:255" json:"location,omitempty"`
	Price       float64 `gorm:"not null" json:"price"`
	VIPPrice    float64 `gorm:"column:vip_price;not null" json:"vip_price"` // Penyesuaian kolom
	Description string  `gorm:"type:text" json:"description"`
	CreatedAt   string  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
