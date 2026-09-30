package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// RequireAPIKey adalah fungsi buat ngecek header sebelum request diproses
func RequireAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")
		deviceID := c.Get("X-DEVICE-ID")

		// Ambil API_KEY murni dari file .env (tanpa fallback)
		expectedKey := os.Getenv("API_KEY")
		
		// Kalau di server (.env) belum diset, tolak semua akses!
		if expectedKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Server Misconfiguration: API_KEY belum diset di environment",
			})
		}

		// Cek validitas API Key dari request client
		if apiKey != expectedKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized: X-API-KEY tidak valid",
			})
		}

		// Cek Device ID
		if deviceID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Bad Request: Header X-DEVICE-ID wajib diisi",
			})
		}

		// Kalau aman semua, lanjut ke Controller (Handler)
		return c.Next()
	}
}