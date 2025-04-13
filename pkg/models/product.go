package models

import (
	"github.com/jinzhu/gorm"
)

// Category ürün kategorisi
type Category struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	ParentID    *uint
	Parent      *Category `gorm:"foreignkey:ParentID"`
	Products    []Product
}

// Unit birim tipi (kg, adet, litre vb.)
type Unit string

const (
	UnitKg    Unit = "kg"
	UnitLiter Unit = "liter"
	UnitPiece Unit = "piece"
	UnitBox   Unit = "box"
)

// Product ürün modeli
type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Barcode     string `gorm:"unique;index"`
	Description string
	CategoryID  uint
	Category    Category `gorm:"foreignkey:CategoryID"`
	Unit        Unit     `gorm:"type:string;default:'piece'"`
	MinStock    float64  `gorm:"default:0"`
	ShelfLife   int      `gorm:"default:0"` // Raf ömrü (gün)
	Price       float64  `gorm:"default:0"`
}

// StockMovement stok hareketleri
type StockMovement struct {
	gorm.Model
	ProductID  uint
	Product    Product `gorm:"foreignkey:ProductID"`
	LocationID uint
	Location   Location `gorm:"foreignkey:LocationID"`
	UserID     uint
	User       User   `gorm:"foreignkey:UserID"`
	Type       string `gorm:"type:string"` // "in" veya "out"
	Quantity   float64
	Note       string
}

// Location depo veya satış lokasyonu
type Location struct {
	gorm.Model
	Name        string `gorm:"not null;unique"`
	Description string
	Address     string
	IsActive    bool `gorm:"default:true"`
}

// Stock stok tablosu (ProductID + LocationID bazında hesaplanır)
type Stock struct {
	gorm.Model
	ProductID  uint
	Product    Product `gorm:"foreignkey:ProductID"`
	LocationID uint
	Location   Location `gorm:"foreignkey:LocationID"`
	Quantity   float64  `gorm:"default:0"`
}
