package domain

// NotificationProvider adalah kontrak yang wajib dipenuhi oleh semua provider (Barantum/SMTP/dll)
type NotificationProvider interface {
	// Fungsi utama buat ngirim OTP
	SendOTP(recipient string, code string) error
	
	// Buat ngasih tau sistem ini provider apa (misal: "BARANTUM")
	GetProviderName() string
	
	// Buat ngasih tau sistem ini lewat jalur apa (misal: "WHATSAPP")
	GetChannelName() string 
}