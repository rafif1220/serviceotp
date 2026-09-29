package domain

import "time"

// OtpLog merepresentasikan tabel otp_logs di MySQL
type OtpLog struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Recipient    string    `json:"recipient"`
	Channel      string    `json:"channel"`  // WHATSAPP atau EMAIL
	Provider     string    `json:"provider"` // BARANTUM atau SMTP
	Status       string    `json:"status"`   // PENDING, SENT, FAILED, VERIFIED, EXPIRED
	DeviceID     string    `json:"device_id"`
	IPAddress    string    `json:"ip_address"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OTPRepository adalah kontrak buat insert/update data log historis ke MySQL
type OTPRepository interface {
	CreateLog(log *OtpLog) error
	UpdateStatus(id string, status string) error
}

// CacheRepository adalah kontrak buat manipulasi data sementara (OTP & Rate Limit) di Redis
type CacheRepository interface {
	SetOTP(recipient string, code string, ttl time.Duration) error
	GetOTP(recipient string) (string, error)
	DeleteOTP(recipient string) error
	IncrementRateLimit(key string, ttl time.Duration) (int, error)
}