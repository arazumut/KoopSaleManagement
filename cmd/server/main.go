package main

import (
	"log"
	"os"

	"koopsatis/pkg/database"
	"koopsatis/pkg/handlers"
	"koopsatis/pkg/middleware"
	"koopsatis/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Println("Uyarı: .env dosyası bulunamadı, varsayılan değerler kullanılacak")
	}

	// Geliştirme modunda Gin'i release moduna al
	if os.Getenv("ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// Veritabanı bağlantısını başlat
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kurulamadı: %v", err)
	}
	defer database.CloseDB()

	// Admin kullanıcı oluştur
	if err := database.CreateAdminUser(db); err != nil {
		log.Printf("Admin kullanıcı oluşturulurken hata: %v", err)
	}

	// Gin router'ı oluştur
	router := gin.Default()

	// CORS ayarları
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API route'ları
	api := router.Group("/api")
	{
		// Auth route'ları (kimlik doğrulama gerektirmeyen)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
		}

		// Kullanıcı yönetimi route'ları
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/", middleware.RoleAuth(models.AdminRole), handlers.GetUsers)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", middleware.RoleAuth(models.AdminRole), handlers.DeleteUser)
		}

		// Ürün yönetimi route'ları
		products := api.Group("/products")
		products.Use(middleware.AuthMiddleware())
		{
			products.POST("/", middleware.RoleAuth(models.AdminRole, models.StockRole), handlers.CreateProduct)
			products.GET("/", handlers.GetProducts)
			products.GET("/:id", handlers.GetProduct)
			products.PUT("/:id", middleware.RoleAuth(models.AdminRole, models.StockRole), handlers.UpdateProduct)
			products.DELETE("/:id", middleware.RoleAuth(models.AdminRole), handlers.DeleteProduct)
		}

		// Diğer API route'ları buraya eklenecek
		// ...

	}

	// Statik dosyaları sunmak için
	router.Static("/static", "./static")

	// Şablon klasörünün yolunu düzelttik
	router.LoadHTMLGlob("/Users/umutaraz/Desktop/KoopSaleManagement/templates/*")

	// Ana sayfa
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Kooperatif Ürün Satış ve Stok Takip Sistemi",
		})
	})

	// Sunucuyu başlat
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Sunucu %s portunda başlatılıyor...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
