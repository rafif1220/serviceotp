package barantum

import (
	"log"

	"github.com/rafif1220/serviceotp/internal/domain"
)

type barantumProvider struct {
	// Nanti bisa diisi config kayak API Key & URL Barantum
}

func NewBarantumProvider() domain.NotificationProvider {
	return &barantumProvider{}
}

func (b *barantumProvider) SendOTP(recipient string, code string) error {
	// TODO: Tulis kodingan HTTP POST ke endpoint Barantum di sini
	// Sementara kita mock pake log print dulu
	log.Printf("[BARANTUM - WHATSAPP] Mengirim OTP %s ke nomor %s", code, recipient)
	return nil
}

func (b *barantumProvider) GetProviderName() string { return "BARANTUM" }
func (b *barantumProvider) GetChannelName() string  { return "WHATSAPP" }