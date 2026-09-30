package smtp

import (
	"log"

	"github.com/rafif1220/serviceotp/internal/domain"
)

type smtpProvider struct {
	// Nanti diisi host, port, email, password dari .env
}

func NewSMTPProvider() domain.NotificationProvider {
	return &smtpProvider{}
}

func (s *smtpProvider) SendOTP(recipient string, code string) error {
	// TODO: Tulis kodingan net/smtp SendMail di sini
	// Sementara kita mock pake log print dulu
	log.Printf("[SMTP - EMAIL] Mengirim OTP %s ke email %s", code, recipient)
	return nil
}

func (s *smtpProvider) GetProviderName() string { return "SMTP" }
func (s *smtpProvider) GetChannelName() string  { return "EMAIL" }