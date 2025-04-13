package handlers

import (
	"net/http"
	"strconv"
	"time"

	"koopsatis/pkg/database"
	"koopsatis/pkg/models"
	"koopsatis/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserLoginRequest login istek modeli
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserCreateRequest kullanıcı oluşturma istek modeli
type UserCreateRequest struct {
	Username  string      `json:"username" binding:"required"`
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password" binding:"required,min=6"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      models.Role `json:"role"`
}

// UserResponse kullanıcı yanıt modeli
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserLoginResponse login yanıt modeli
type UserLoginResponse struct {
	Token    string       `json:"token"`
	User     UserResponse `json:"user"`
	ExpireAt time.Time    `json:"expire_at"`
}

// Login kullanıcı giriş işlemi
func Login(c *gin.Context) {
	var request UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kullanıcıyı veritabanında ara
	var user models.User
	if err := database.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "geçersiz kullanıcı adı veya şifre"})
		return
	}

	// Şifreyi kontrol et
	if !user.CheckPassword(request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "geçersiz kullanıcı adı veya şifre"})
		return
	}

	// Token oluştur
	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token oluşturulamadı"})
		return
	}

	// Kullanıcı son giriş zamanını güncelle
	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	// Aktivite logu kaydet
	logActivity := models.ActivityLog{
		UserID:     user.ID,
		IP:         c.ClientIP(),
		ActionType: models.ActionLogin,
		EntityType: models.EntityUser,
		EntityID:   &user.ID,
		Details:    "Kullanıcı giriş yaptı",
		Timestamp:  time.Now(),
	}
	database.DB.Create(&logActivity)

	// Yanıt hazırla
	response := UserLoginResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	c.JSON(http.StatusOK, response)
}

// Register kullanıcı kayıt işlemi
func Register(c *gin.Context) {
	var request UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Admin değilse, sadece üye rolünde kayıt olabilsin
	userRole, exists := c.Get("role")
	isAdmin := exists && userRole.(string) == string(models.AdminRole)

	if !isAdmin && request.Role != "" && request.Role != models.MemberRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "sadece üye rolü ile kayıt olabilirsiniz"})
		return
	}

	// Kullanıcı adı ve email kontrolü
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", request.Username, request.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bu kullanıcı adı veya e-posta adresi zaten kullanılıyor"})
		return
	}

	// Yeni kullanıcı oluştur
	user := models.User{
		Username:  request.Username,
		Email:     request.Email,
		Password:  request.Password, // BeforeSave hook'u ile hashlenecek
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Role:      request.Role,
		Active:    true,
	}

	// Eğer rol belirtilmemişse, varsayılan olarak üye rolü ver
	if user.Role == "" {
		user.Role = models.MemberRole
	}

	// Kullanıcıyı kaydet
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kullanıcı oluşturulamadı: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// GetUsers tüm kullanıcıları listeler (sadece admin)
func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetUser belirli bir kullanıcıyı getirir
func GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz kullanıcı kimliği"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kullanıcı bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// UpdateUser kullanıcı bilgilerini günceller
func UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz kullanıcı kimliği"})
		return
	}

	// Mevcut kullanıcıyı bul
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kullanıcı bulunamadı"})
		return
	}

	// İsteği bağla
	var request struct {
		Email     *string      `json:"email"`
		Password  *string      `json:"password"`
		FirstName *string      `json:"first_name"`
		LastName  *string      `json:"last_name"`
		Role      *models.Role `json:"role"`
		Active    *bool        `json:"active"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Yetki kontrolü - sadece admin diğer kullanıcıların rolünü değiştirebilir
	currentUserID, _ := c.Get("userID")
	currentUserRole, _ := c.Get("role")
	isAdmin := currentUserRole.(string) == string(models.AdminRole)
	isSelf := currentUserID.(uint) == uint(userID)

	// Kullanıcı sadece kendi bilgilerini güncelleyebilir veya admin tüm bilgileri
	if !isAdmin && !isSelf {
		c.JSON(http.StatusForbidden, gin.H{"error": "bu işlem için yetkiniz yok"})
		return
	}

	// Admin olmayan kullanıcılar rol değiştiremez
	if !isAdmin && request.Role != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "rol değiştirme yetkiniz yok"})
		return
	}

	// Güncelleme işlemi
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Password != nil {
		user.Password = *request.Password
	}
	if request.FirstName != nil {
		user.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		user.LastName = *request.LastName
	}
	if request.Role != nil && isAdmin {
		user.Role = *request.Role
	}
	if request.Active != nil && isAdmin {
		user.Active = *request.Active
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kullanıcı güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// DeleteUser kullanıcıyı siler (sadece admin)
func DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz kullanıcı kimliği"})
		return
	}

	// Kullanıcıyı bul
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kullanıcı bulunamadı"})
		return
	}

	// Admin kullanıcısını silmeye izin verme
	if user.Role == models.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin kullanıcısı silinemez"})
		return
	}

	// Kullanıcıyı sil (soft delete)
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kullanıcı silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kullanıcı başarıyla silindi"})
}
