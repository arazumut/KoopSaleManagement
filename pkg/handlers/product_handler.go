package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"koopsatis/pkg/database"
	"koopsatis/pkg/models"

	"github.com/gin-gonic/gin"
)

// ProductRequest ürün oluşturma/güncelleme için istek modeli
type ProductRequest struct {
	Name        string      `json:"name" binding:"required"`
	Barcode     string      `json:"barcode"`
	Description string      `json:"description"`
	CategoryID  uint        `json:"category_id" binding:"required"`
	Unit        models.Unit `json:"unit"`
	MinStock    float64     `json:"min_stock"`
	ShelfLife   int         `json:"shelf_life"`
	Price       float64     `json:"price" binding:"required"`
}

// ProductResponse ürün yanıt modeli
type ProductResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Barcode     string      `json:"barcode"`
	Description string      `json:"description"`
	CategoryID  uint        `json:"category_id"`
	Category    string      `json:"category"`
	Unit        string      `json:"unit"`
	MinStock    float64     `json:"min_stock"`
	ShelfLife   int         `json:"shelf_life"`
	Price       float64     `json:"price"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	StockInfo   []StockInfo `json:"stock_info,omitempty"`
}

// StockInfo stok bilgisi yanıt modeli
type StockInfo struct {
	LocationID   uint    `json:"location_id"`
	LocationName string  `json:"location_name"`
	Quantity     float64 `json:"quantity"`
}

// GenerateBarcode rastgele barkod oluşturur
func GenerateBarcode() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	barcode := "978"
	for i := 0; i < 10; i++ {
		barcode += string(digits[rand.Intn(len(digits))])
	}
	return barcode
}

// CreateProduct yeni ürün oluşturur
func CreateProduct(c *gin.Context) {
	var request ProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kategori kontrolü
	var category models.Category
	if err := database.DB.First(&category, request.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kategori bulunamadı"})
		return
	}

	// Barkod otomatik oluştur
	barcode := request.Barcode
	if barcode == "" {
		barcode = GenerateBarcode()
	} else {
		// Mevcut barkod kontrolü
		var existingProduct models.Product
		if err := database.DB.Where("barcode = ?", barcode).First(&existingProduct).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bu barkod zaten kullanılıyor"})
			return
		}
	}

	// Birim değeri kontrolü
	unit := request.Unit
	if unit == "" {
		unit = models.UnitPiece
	}

	// Yeni ürün oluştur
	product := models.Product{
		Name:        request.Name,
		Barcode:     barcode,
		Description: request.Description,
		CategoryID:  request.CategoryID,
		Unit:        unit,
		MinStock:    request.MinStock,
		ShelfLife:   request.ShelfLife,
		Price:       request.Price,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ürün oluşturulamadı: " + err.Error()})
		return
	}

	// Kullanıcı ID'sini al ve log kaydı ekle
	userID, exists := c.Get("userID")
	if exists {
		logActivity := models.ActivityLog{
			UserID:     userID.(uint),
			IP:         c.ClientIP(),
			ActionType: models.ActionCreate,
			EntityType: models.EntityProduct,
			EntityID:   &product.ID,
			Details:    "Yeni ürün oluşturuldu: " + product.Name,
			Timestamp:  time.Now(),
		}
		database.DB.Create(&logActivity)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ürün başarıyla oluşturuldu",
		"product": ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Barcode:     product.Barcode,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			Category:    category.Name,
			Unit:        string(product.Unit),
			MinStock:    product.MinStock,
			ShelfLife:   product.ShelfLife,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	})
}

// GetProducts tüm ürünleri listeler
func GetProducts(c *gin.Context) {
	// Filtre parametreleri
	categoryID := c.Query("category_id")
	name := c.Query("name")
	barcode := c.Query("barcode")

	var products []models.Product
	query := database.DB.Model(&models.Product{}).Preload("Category")

	// Filtreleri uygula
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if barcode != "" {
		query = query.Where("barcode = ?", barcode)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ürünler getirilemedi"})
		return
	}

	// Yanıt oluştur
	var response []ProductResponse
	for _, product := range products {
		response = append(response, ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Barcode:     product.Barcode,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			Category:    product.Category.Name,
			Unit:        string(product.Unit),
			MinStock:    product.MinStock,
			ShelfLife:   product.ShelfLife,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct belirli bir ürünü getirir
func GetProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz ürün kimliği"})
		return
	}

	var product models.Product
	if err := database.DB.Preload("Category").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ürün bulunamadı"})
		return
	}

	// Stok bilgisini getir
	var stocks []models.Stock
	database.DB.Where("product_id = ?", productID).Preload("Location").Find(&stocks)

	var stockInfo []StockInfo
	for _, stock := range stocks {
		stockInfo = append(stockInfo, StockInfo{
			LocationID:   stock.LocationID,
			LocationName: stock.Location.Name,
			Quantity:     stock.Quantity,
		})
	}

	c.JSON(http.StatusOK, ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Barcode:     product.Barcode,
		Description: product.Description,
		CategoryID:  product.CategoryID,
		Category:    product.Category.Name,
		Unit:        string(product.Unit),
		MinStock:    product.MinStock,
		ShelfLife:   product.ShelfLife,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		StockInfo:   stockInfo,
	})
}

// UpdateProduct ürün bilgilerini günceller
func UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz ürün kimliği"})
		return
	}

	// Mevcut ürünü bul
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ürün bulunamadı"})
		return
	}

	var request ProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kategori kontrolü
	var category models.Category
	if err := database.DB.First(&category, request.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kategori bulunamadı"})
		return
	}

	// Barkod kontrolü (yeni barkod verilmişse ve başka ürün tarafından kullanılıyorsa)
	if request.Barcode != "" && request.Barcode != product.Barcode {
		var existingProduct models.Product
		if err := database.DB.Where("barcode = ? AND id != ?", request.Barcode, productID).First(&existingProduct).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bu barkod zaten başka bir ürün tarafından kullanılıyor"})
			return
		}
		product.Barcode = request.Barcode
	}

	// Ürünü güncelle
	product.Name = request.Name
	product.Description = request.Description
	product.CategoryID = request.CategoryID
	if request.Unit != "" {
		product.Unit = request.Unit
	}
	product.MinStock = request.MinStock
	product.ShelfLife = request.ShelfLife
	product.Price = request.Price

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ürün güncellenemedi"})
		return
	}

	// Kullanıcı ID'sini al ve log kaydı ekle
	userID, exists := c.Get("userID")
	if exists {
		logActivity := models.ActivityLog{
			UserID:     userID.(uint),
			IP:         c.ClientIP(),
			ActionType: models.ActionUpdate,
			EntityType: models.EntityProduct,
			EntityID:   &product.ID,
			Details:    "Ürün güncellendi: " + product.Name,
			Timestamp:  time.Now(),
		}
		database.DB.Create(&logActivity)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ürün başarıyla güncellendi",
		"product": ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Barcode:     product.Barcode,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			Category:    category.Name,
			Unit:        string(product.Unit),
			MinStock:    product.MinStock,
			ShelfLife:   product.ShelfLife,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	})
}

// DeleteProduct ürünü siler
func DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz ürün kimliği"})
		return
	}

	// Ürünü bul
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ürün bulunamadı"})
		return
	}

	// Stok kontrolü (ürünün stokta olup olmadığını kontrol et)
	var stockCount int64
	database.DB.Model(&models.Stock{}).Where("product_id = ? AND quantity > 0", productID).Count(&stockCount)
	if stockCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stokta bulunan ürün silinemez, önce stoğu sıfırlayın"})
		return
	}

	// Satışlarda kullanılıp kullanılmadığının kontrolü
	var saleItemCount int64
	database.DB.Model(&models.SaleItem{}).Where("product_id = ?", productID).Count(&saleItemCount)
	if saleItemCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "satışlarda kullanılan ürün silinemez"})
		return
	}

	// Ürünü sil (soft delete)
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ürün silinemedi"})
		return
	}

	// Kullanıcı ID'sini al ve log kaydı ekle
	userID, exists := c.Get("userID")
	if exists {
		logActivity := models.ActivityLog{
			UserID:     userID.(uint),
			IP:         c.ClientIP(),
			ActionType: models.ActionDelete,
			EntityType: models.EntityProduct,
			EntityID:   &product.ID,
			Details:    "Ürün silindi: " + product.Name,
			Timestamp:  time.Now(),
		}
		database.DB.Create(&logActivity)
	}

	c.JSON(http.StatusOK, gin.H{"message": "ürün başarıyla silindi"})
}
