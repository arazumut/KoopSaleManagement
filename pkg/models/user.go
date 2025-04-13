package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Role kullanıcı rolleri için enum tipi
type Role string

const (
	AdminRole  Role = "admin"
	MemberRole Role = "member"
	SalesRole  Role = "sales"
	StockRole  Role = "stock"
)

// User kullanıcı modeli
type User struct {
	gorm.Model
	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	FirstName   string
	LastName    string
	Role        Role `gorm:"type:string;default:'member'"`
	Active      bool `gorm:"default:true"`
	LastLoginAt *time.Time
}

// HashPassword şifreyi bcrypt ile hashler
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verilen şifre ile hashlenen şifreyi karşılaştırır
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeSave GORM hook'u ile kaydetmeden önce şifreyi hashler
func (u *User) BeforeSave() error {
	// Eğer şifre hashli değilse
	if len(u.Password) > 0 && len(u.Password) < 60 {
		return u.HashPassword()
	}
	return nil
}
