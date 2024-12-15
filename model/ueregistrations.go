package model

// Struktur untuk registrasi pengguna pada acara (User_Event_Registrations)
type UserEventRegistration struct {
	RegistrationID   uint   `gorm:"primaryKey;autoIncrement" json:"registration_id"`
	UserID           uint   `gorm:"not null" json:"user_id"`
	EventID          uint   `gorm:"not null" json:"event_id"`
	Status           string `gorm:"default:'biasa'" json:"status"`                      // 'biasa', 'vip', 'reject'
	PaymentReceipt   string `gorm:"size:255" json:"payment_receipt"`                    // Lokasi bukti pembayaran
	MateriFile       string `gorm:"size:255" json:"materi_file,omitempty"`              // Lokasi file materi untuk user VIP
	SertifikatFile   string `gorm:"size:255" json:"sertifikat_file,omitempty"`          // Lokasi file sertifikat untuk user VIP
	RegistrationDate string `gorm:"default:CURRENT_TIMESTAMP" json:"registration_date"` // Tanggal registrasi
	User             User   `gorm:"foreignKey:UserID" json:"user"`
	Event            Event  `gorm:"foreignKey:EventID" json:"event"`
}
