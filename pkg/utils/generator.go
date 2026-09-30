package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateOTP bikin kode angka random sesuai panjang yang diminta (misal: 6 digit)
func GenerateOTP(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}