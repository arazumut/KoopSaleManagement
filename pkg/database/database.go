package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"koopsatis/pkg/models"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var DB *gorm.DB

// InitDB veritabanı bağlantısını başlatır
func InitDB() (*gorm.DB, error) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbPath := os.Getenv("DB_PATH")

	var err error
	if dbDriver == "sqlite3" {
		DB, err = gorm.Open(dbDriver, dbPath)
	} else {
		return nil, fmt.Errorf("desteklenmeyen veritabanı tipi: %s", dbDriver)
	}

	if err != nil {
		return nil, err
	}

	// SQLite özellik ayarları
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(time.Hour)

	// Hata ayıklama modunda SQL sorgularını göster (geliştirme sırasında)
	if os.Getenv("ENV") == "development" {
		DB.LogMode(true)
	}

	// Tabloları oluştur (Migrasyon)
	err = RunMigrations(DB)
	if err != nil {
		return nil, err
	}

	return DB, nil
}

// RunMigrations tabloları oluşturur
func RunMigrations(db *gorm.DB) error {
	log.Println("Veritabanı tabloları oluşturuluyor...")

	// Modelleri otomatik migrate edelim
	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Location{},
		&models.Stock{},
		&models.StockMovement{},
		&models.Customer{},
		&models.Sale{},
		&models.SaleItem{},
		&models.Payment{},
		&models.SaleReturn{},
		&models.SaleReturnItem{},
		&models.Supplier{},
		&models.Purchase{},
		&models.PurchaseItem{},
		&models.SupplierPayment{},
		&models.ActivityLog{},
		&models.SystemLog{},
	).Error

	if err != nil {
		return err
	}

	// Foreign key ilişkilerini oluştur
	// ...

	log.Println("Veritabanı tabloları başarıyla oluşturuldu.")
	return nil
}

// CreateAdminUser eğer admin kullanıcı yoksa oluşturur
func CreateAdminUser(db *gorm.DB) error {
	var adminCount int
	db.Model(&models.User{}).Where("role = ?", models.AdminRole).Count(&adminCount)

	if adminCount == 0 {
		adminUser := models.User{
			Username:  "admin",
			Email:     "admin@example.com",
			Password:  "admin123", // otomatik hashlenecek
			FirstName: "Admin",
			LastName:  "User",
			Role:      models.AdminRole,
			Active:    true,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Println("Admin kullanıcı başarıyla oluşturuldu.")
	}

	return nil
}

// CloseDB veritabanı bağlantısını kapatır
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
