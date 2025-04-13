package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// ActionType eylem tipini tanımlar
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionView   ActionType = "view"
	ActionLogin  ActionType = "login"
	ActionLogout ActionType = "logout"
)

// EntityType etkilenen veri tipini tanımlar
type EntityType string

const (
	EntityUser     EntityType = "user"
	EntityProduct  EntityType = "product"
	EntityCategory EntityType = "category"
	EntityCustomer EntityType = "customer"
	EntitySupplier EntityType = "supplier"
	EntitySale     EntityType = "sale"
	EntityPurchase EntityType = "purchase"
	EntityStock    EntityType = "stock"
	EntitySystem   EntityType = "system"
)

// ActivityLog kullanıcı eylem kayıtları
type ActivityLog struct {
	gorm.Model
	UserID     uint
	User       User `gorm:"foreignkey:UserID"`
	IP         string
	ActionType ActionType `gorm:"type:string"`
	EntityType EntityType `gorm:"type:string"`
	EntityID   *uint
	Details    string
	Timestamp  time.Time
}

// SystemLog sistem olay kayıtları
type SystemLog struct {
	ID        uint   `gorm:"primarykey"`
	Level     string `gorm:"type:string;default:'info'"`
	Message   string
	Source    string
	Timestamp time.Time `gorm:"index"`
}
