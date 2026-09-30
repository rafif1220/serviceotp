package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rafif1220/serviceotp/internal/domain"
	"github.com/rafif1220/serviceotp/pkg/utils"
)

// OtpUseCase adalah kontrak untuk layer HTTP nanti
type OtpUseCase interface {
	RequestOTP(recipient string, deviceID string, ipAddress string) error
	// Nanti kita bisa nambahin VerifyOTP di sini
}

type otpUseCaseImpl struct {
	mysqlRepo domain.OTPRepository
	redisRepo domain.CacheRepository
	barantum  domain.NotificationProvider
	smtp      domain.NotificationProvider
}

// NewOtpUseCase adalah constructor buat nyuntik semua dependency
func NewOtpUseCase(
	mysql domain.OTPRepository,
	redis domain.CacheRepository,
	barantum domain.NotificationProvider,
	smtp domain.NotificationProvider,
) OtpUseCase {
	return &otpUseCaseImpl{
		mysqlRepo: mysql,
		redisRepo: redis,
		barantum:  barantum,
		smtp:      smtp,
	}
}

func (u *otpUseCaseImpl) RequestOTP(recipient string, deviceID string, ipAddress string) error {
	// 1. Cek Rate Limit (Maksimal 3x request per jam untuk nomor/email yang sama)
	rateLimitKey := "nukar:ratelimit:target:" + recipient
	count, err := u.redisRepo.IncrementRateLimit(rateLimitKey, 1*time.Hour)
	if err != nil {
		return err
	}
	if count > 3 {
		return errors.New("terlalu banyak request, coba lagi dalam 1 jam")
	}

	// 2. Tentukan Provider berdasarkan format recipient (Strategy Pattern)
	var provider domain.NotificationProvider
	if strings.Contains(recipient, "@") {
		provider = u.smtp
	} else {
		provider = u.barantum
	}

	// 3. Generate Kode OTP (6 digit)
	otpCode, err := utils.GenerateOTP(6)
	if err != nil {
		return err
	}

	// 4. Simpan OTP ke Redis (TTL 5 Menit)
	err = u.redisRepo.SetOTP(recipient, otpCode, 5*time.Minute)
	if err != nil {
		return err
	}

	// 5. Catat log PENDING ke MySQL
	logID := uuid.New().String()
	otpLog := &domain.OtpLog{
		ID:        logID,
		Recipient: recipient,
		Channel:   provider.GetChannelName(),
		Provider:  provider.GetProviderName(),
		Status:    "PENDING",
		DeviceID:  deviceID,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	u.mysqlRepo.CreateLog(otpLog)

	// 6. Eksekusi Pengiriman OTP via Provider terpilih
	err = provider.SendOTP(recipient, otpCode)
	
	// 7. Update status log di MySQL berdasarkan hasil pengiriman
	if err != nil {
		u.mysqlRepo.UpdateStatus(logID, "FAILED")
		return errors.New("gagal mengirim OTP: " + err.Error())
	}

	u.mysqlRepo.UpdateStatus(logID, "SENT")
	return nil
}