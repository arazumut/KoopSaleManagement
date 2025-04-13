package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Supplier tedarikçi modeli
type Supplier struct {
	gorm.Model
	Name        string `gorm:"not null"`
	ContactName string
	PhoneNumber string
	Email       string
	Address     string
	TaxNumber   string
	TaxOffice   string
	Notes       string
	Purchases   []Purchase
	Balance     float64 `gorm:"default:0"` // Bakiye/borç durumu (- alacak, + borç)
}

// Purchase satın alma modeli
type Purchase struct {
	gorm.Model
	InvoiceNumber string `gorm:"unique;index"`
	SupplierID    uint
	Supplier      Supplier `gorm:"foreignkey:SupplierID"`
	UserID        uint
	User          User `gorm:"foreignkey:UserID"`
	PurchaseDate  time.Time
	DueDate       *time.Time
	TotalAmount   float64
	TaxAmount     float64
	Discount      float64
	FinalAmount   float64
	IsPaid        bool
	Note          string
	Items         []PurchaseItem
	Payments      []SupplierPayment
}

// PurchaseItem satın alma kalemi
type PurchaseItem struct {
	gorm.Model
	PurchaseID uint
	Purchase   Purchase `gorm:"foreignkey:PurchaseID"`
	ProductID  uint
	Product    Product `gorm:"foreignkey:ProductID"`
	Quantity   float64
	UnitPrice  float64
	Discount   float64
	TaxRate    float64 `gorm:"default:0.18"` // Türkiye'de standart KDV oranı
	TaxAmount  float64
	Subtotal   float64
	Total      float64
}

// SupplierPayment tedarikçi ödeme modeli
type SupplierPayment struct {
	gorm.Model
	PurchaseID      *uint
	Purchase        *Purchase `gorm:"foreignkey:PurchaseID"`
	SupplierID      uint
	Supplier        Supplier `gorm:"foreignkey:SupplierID"`
	Amount          float64
	PaymentType     PaymentType `gorm:"type:string;default:'cash'"`
	PaymentDate     time.Time
	ReferenceNo     string // Banka EFT numarası, çek numarası vb.
	Note            string
	ProcessedByID   uint
	ProcessedByUser User `gorm:"foreignkey:ProcessedByID"`
}
