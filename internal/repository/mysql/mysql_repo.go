package mysqlrepo

import (
	"github.com/rafif1220/serviceotp/internal/domain"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

// NewOTPRepository constructor buat repo MySQL
func NewOTPRepository(db *gorm.DB) domain.OTPRepository {
	return &mysqlRepo{db: db}
}

func (m *mysqlRepo) CreateLog(log *domain.OtpLog) error {
	return m.db.Create(log).Error
}

func (m *mysqlRepo) UpdateStatus(id string, status string) error {
	// Cuma nge-update kolom "status" berdasarkan ID
	return m.db.Model(&domain.OtpLog{}).Where("id = ?", id).Update("status", status).Error
}
