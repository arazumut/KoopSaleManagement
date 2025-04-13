package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// PaymentType ödeme tipi
type PaymentType string

const (
	PaymentCash   PaymentType = "cash"
	PaymentCard   PaymentType = "card"
	PaymentCredit PaymentType = "credit"
	PaymentBank   PaymentType = "bank_transfer"
)

// Customer müşteri modeli
type Customer struct {
	gorm.Model
	Name        string `gorm:"not null"`
	PhoneNumber string
	Email       string
	Address     string
	TaxNumber   string
	TaxOffice   string
	Notes       string
	Sales       []Sale
	Payments    []Payment
	Balance     float64 `gorm:"default:0"` // Bakiye/borç durumu (- borç, + alacak)
}

// Sale satış modeli
type Sale struct {
	gorm.Model
	InvoiceNumber string `gorm:"unique;index"`
	CustomerID    uint
	Customer      Customer `gorm:"foreignkey:CustomerID"`
	UserID        uint
	User          User `gorm:"foreignkey:UserID"`
	SaleDate      time.Time
	DueDate       *time.Time
	TotalAmount   float64
	TaxAmount     float64
	Discount      float64
	FinalAmount   float64
	IsPaid        bool
	Note          string
	Items         []SaleItem
	Payments      []Payment
}

// SaleItem satış kalemi
type SaleItem struct {
	gorm.Model
	SaleID    uint
	Sale      Sale `gorm:"foreignkey:SaleID"`
	ProductID uint
	Product   Product `gorm:"foreignkey:ProductID"`
	Quantity  float64
	UnitPrice float64
	Discount  float64
	TaxRate   float64 `gorm:"default:0.18"` // Türkiye'de standart KDV oranı
	TaxAmount float64
	Subtotal  float64
	Total     float64
}

// Payment ödeme modeli
type Payment struct {
	gorm.Model
	SaleID         *uint
	Sale           *Sale `gorm:"foreignkey:SaleID"`
	CustomerID     uint
	Customer       Customer `gorm:"foreignkey:CustomerID"`
	Amount         float64
	PaymentType    PaymentType `gorm:"type:string;default:'cash'"`
	PaymentDate    time.Time
	ReferenceNo    string // Banka EFT numarası, çek numarası vb.
	Note           string
	ReceivedByID   uint
	ReceivedByUser User `gorm:"foreignkey:ReceivedByID"`
}

// SaleReturnItem iade edilen ürün kalemi
type SaleReturnItem struct {
	gorm.Model
	SaleReturnID uint
	SaleReturn   SaleReturn `gorm:"foreignkey:SaleReturnID"`
	SaleItemID   uint
	SaleItem     SaleItem `gorm:"foreignkey:SaleItemID"`
	ProductID    uint
	Product      Product `gorm:"foreignkey:ProductID"`
	Quantity     float64
	UnitPrice    float64
	TaxRate      float64
	TaxAmount    float64
	Total        float64
	Reason       string
}

// SaleReturn iade modeli
type SaleReturn struct {
	gorm.Model
	SaleID      uint
	Sale        Sale `gorm:"foreignkey:SaleID"`
	ReturnDate  time.Time
	TotalAmount float64
	Reason      string
	ProcessedBy uint
	User        User `gorm:"foreignkey:ProcessedBy"`
	Items       []SaleReturnItem
}
