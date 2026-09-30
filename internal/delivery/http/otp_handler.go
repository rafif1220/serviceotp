package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafif1220/serviceotp/internal/usecase"
)

type OtpHandler struct {
	usecase usecase.OtpUseCase
}

// NewOtpHandler adalah constructor buat Controller OTP
func NewOtpHandler(u usecase.OtpUseCase) *OtpHandler {
	return &OtpHandler{usecase: u}
}

// Route buat daftarin endpoint ke Fiber
func (h *OtpHandler) Route(app fiber.Router) {
	// Endpoint: POST /api/v1/otp/request
	app.Post("/request", h.RequestOTP)
}

func (h *OtpHandler) RequestOTP(c *fiber.Ctx) error {
	// Struct buat nangkep JSON dari body request
	var req struct {
		Recipient string `json:"recipient"`
	}

	// Parsing body JSON ke struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request JSON tidak valid",
		})
	}

	// Ambil metadata dari request
	deviceID := c.Get("X-DEVICE-ID")
	ipAddress := c.IP()

	// Lempar ke layer Usecase buat diproses (cek redis, simpan mysql, kirim wa/email)
	err := h.usecase.RequestOTP(req.Recipient, deviceID, ipAddress)
	if err != nil {
		// Uncomment baris di bawah kalau butuh pantau log error/rate-limit di terminal server (tambahkan import "log" di atas)
		// log.Printf("[WARNING] Request OTP ditolak untuk %s: %v", req.Recipient, err)

		// Balikin status 429 ke client (Postman/Web)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Kalau sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "OTP berhasil diproses dan sedang dikirim",
	})
}