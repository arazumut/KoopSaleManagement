package main

import (
	"html/template"
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

	// Template fonksiyonları ekle
	router.SetFuncMap(template.FuncMap{
		"formatDate":     handlers.FormatDate,
		"formatCurrency": handlers.FormatCurrency,
		"safeHTML":       handlers.SafeHTML,
		"add":            handlers.Add,
		"subtract":       handlers.Subtract,
		"multiply":       handlers.Multiply,
		"divide":         handlers.Divide,
		"hasRole":        handlers.HasRole,
		"error":          func(err interface{}) interface{} { return err },
	})

	// Şablon klasörünün yolunu belirt
	router.LoadHTMLGlob("templates/*")

	// Ana sayfa
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "layout.html", gin.H{
			"title": "Kontrol Paneli",
			"user": gin.H{
				"name": "Admin Kullanıcı",
				"role": "admin",
			},
		})
	})

	// Giriş sayfası route'u
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{
			"title": "Giriş",
		})
	})

	// Ürünler sayfası
	router.GET("/products", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(200, "layout.html", gin.H{
			"title": "Ürünler",
			"user":  c.MustGet("user"),
		})
	})

	// Yeni ürün sayfası
	router.GET("/products/new", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(200, "layout.html", gin.H{
			"title": "Yeni Ürün",
			"user":  c.MustGet("user"),
		})
	})

	// Ürün düzenleme sayfası
	router.GET("/products/:id/edit", middleware.AuthMiddleware(), func(c *gin.Context) {
		productID := c.Param("id")
		// Gerçek uygulamada ürün verisini veritabanından çekeceksiniz
		// Şimdilik örnek veri kullanıyoruz
		c.HTML(200, "layout.html", gin.H{
			"title": "Ürün Düzenle",
			"user":  c.MustGet("user"),
			"Product": gin.H{
				"ID":            productID,
				"Code":          "PRD-001",
				"Name":          "Elma",
				"Category":      "meyve",
				"Unit":          "kg",
				"Stock":         25.0,
				"CriticalStock": 10.0,
				"PurchasePrice": 12.50,
				"SalePrice":     18.00,
				"Description":   "Taze elma",
				"Status":        "active",
			},
		})
	})

	// Satışlar sayfası
	router.GET("/sales", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(200, "layout.html", gin.H{
			"title": "Satışlar",
			"user":  c.MustGet("user"),
		})
	})

	// Yeni satış sayfası
	router.GET("/sales/new", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(200, "layout.html", gin.H{
			"title": "Yeni Satış",
			"user":  c.MustGet("user"),
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
