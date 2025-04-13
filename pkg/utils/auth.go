package utils

import (
	"errors"
	"os"
	"time"

	"koopsatis/pkg/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("geçersiz token")
	ErrExpiredToken = errors.New("token süresi dolmuş")
)

// TokenClaims JWT token içeriği
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken kullanıcı için JWT token oluşturur
func GenerateToken(user *models.User) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	expirationTime := time.Now().Add(24 * time.Hour) // Token 24 saat geçerli

	claims := &TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "koopsatis",
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken JWT tokenin geçerliliğini kontrol eder
func ValidateToken(tokenString string) (*TokenClaims, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
