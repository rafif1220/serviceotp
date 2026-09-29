package config

import (
	"fmt"
	"log"
	"os"

	"github.com/rafif1220/serviceotp/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke MySQL: %v", err)
	}

	// Auto migrate tabel OtpLog biar kita nggak usah CREATE TABLE manual
	err = db.AutoMigrate(&domain.OtpLog{})
	if err != nil {
		log.Fatalf("Gagal migrate tabel: %v", err)
	}

	log.Println("MySQL Connected & Migrated!")
	return db
}