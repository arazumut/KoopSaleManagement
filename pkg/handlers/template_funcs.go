package handlers

import (
	"fmt"
	"html/template"
	"time"

	"koopsatis/pkg/models"
)

// Template yardımcı fonksiyonları

// FormatDate tarihi istenen formatta formatlar
func FormatDate(t time.Time) string {
	return t.Format("02.01.2006")
}

// FormatCurrency para değerini ₺ formatında verir
func FormatCurrency(value float64) string {
	return fmt.Sprintf("%.2f ₺", value)
}

// SafeHTML HTML içeriğini güvenli bir şekilde render eder
func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

// Add toplama işlemi yapar
func Add(a, b int) int {
	return a + b
}

// Subtract çıkarma işlemi yapar
func Subtract(a, b int) int {
	return a - b
}

// Multiply çarpma işlemi yapar
func Multiply(a, b int) int {
	return a * b
}

// Divide bölme işlemi yapar
func Divide(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

// HasRole kullanıcının belirli bir role sahip olup olmadığını kontrol eder
func HasRole(user interface{}, roles ...string) bool {
	if user == nil {
		return false
	}

	// Kullanıcı türüne göre rol kontrolü
	switch u := user.(type) {
	case models.User:
		userRole := string(u.Role)
		for _, role := range roles {
			if userRole == role {
				return true
			}
		}
	case map[string]interface{}:
		if roleVal, ok := u["role"].(string); ok {
			for _, role := range roles {
				if roleVal == role {
					return true
				}
			}
		}
	case map[string]string:
		if roleVal, ok := u["role"]; ok {
			for _, role := range roles {
				if roleVal == role {
					return true
				}
			}
		}
	}

	return false
}
