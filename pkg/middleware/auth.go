package middleware

import (
	"net/http"
	"strings"

	"koopsatis/pkg/models"
	"koopsatis/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware kimlik doğrulama ve yetkilendirme middleware'i
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "yetkilendirme başlığı bulunamadı"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "geçersiz token formatı"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Kullanıcı bilgilerini context'e ekle
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleAuth belirli rollere sahip kullanıcılar için yetkilendirme kontrolü
func RoleAuth(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "kullanıcı rolü bulunamadı"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		authorized := false

		for _, role := range roles {
			if string(role) == roleStr {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "bu işlem için yetkiniz yok"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LogActivity kullanıcı aktivitelerini kaydetme middleware'i
func LogActivity(actionType models.ActionType, entityType models.EntityType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// İstek işlendikten sonra log kaydı ekleyin
		c.Next()

		// Kullanıcı bilgilerini alın
		_, exists := c.Get("userID")
		if !exists {
			return // Giriş yapmamış kullanıcı, loglama yapılmaz
		}

		// Log kaydı ekleme işlemi burada yapılabilir
		// database.DB.Create(&models.ActivityLog{...})
	}
}
