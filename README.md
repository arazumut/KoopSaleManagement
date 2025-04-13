# Kooperatif Ürün Satış ve Stok Takip Sistemi

Kooperatifler için geliştirilmiş, ürün satış ve stok takibi yapabilen web tabanlı bir yönetim sistemidir.

## Özellikler

- **Kullanıcı Yönetimi**: Rol tabanlı yetkilendirme sistemi (Admin, Üye, Satış, Depocu)
- **Ürün Yönetimi**: Ürün listeleme, ekleme, güncelleme, silme
- **Stok Takibi**: Giriş-çıkış işlemleri, lokasyon bazlı stok yönetimi
- **Satış Yönetimi**: Satış formu, faturalandırma, iade işlemleri
- **Müşteri Yönetimi**: Müşteri kaydı, borç ve ödeme takibi
- **Tedarikçi Yönetimi**: Tedarikçi kaydı, alım kayıtları
- **Raporlama**: Satış, stok ve finansal raporlar
- **Kullanıcı Dostu Arayüz**: Responsive tasarım, modern UI

## Teknoloji

- **Backend**: Go (Golang) ve Gin Framework
- **Frontend**: HTML5, TailwindCSS, Alpine.js
- **Veritabanı**: SQLite (PostgreSQL'e geçiş yapılabilir)
- **Kimlik Doğrulama**: JWT token tabanlı

## Kurulum

### Gereksinimler

- Go 1.16 veya daha yüksek sürüm
- SQLite

### Adımlar

1. Depoyu klonlayın:
```bash
git clone https://github.com/kullanici/koopsatis.git
cd koopsatis
```

2. Gerekli bağımlılıkları yükleyin:
```bash
go mod download
```

3. Uygulamayı derleyin:
```bash
go build -o koopsatis ./cmd/server
```

4. Uygulamayı çalıştırın:
```bash
./koopsatis
```

5. Tarayıcınızda aşağıdaki adrese gidin:
```
http://localhost:8080
```

## Geliştirme Ortamı

Geliştirme için aşağıdaki komutu kullanabilirsiniz:

```bash
go run ./cmd/server/main.go
```

## Varsayılan Kullanıcı

Uygulama ilk çalıştırıldığında otomatik olarak bir admin kullanıcısı oluşturulur:

- **Kullanıcı Adı**: admin
- **Şifre**: admin123

İlk girişten sonra güvenlik için bu şifreyi değiştirmeniz önerilir.

## API Dökümantasyonu

RESTful API'lar aşağıdaki base URL üzerinden erişilebilir:

```
http://localhost:8080/api
```

### Kimlik Doğrulama Endpoint'leri

- `POST /api/auth/login`: Kullanıcı girişi
- `POST /api/auth/register`: Yeni kullanıcı kaydı

### Kullanıcı Endpoint'leri

- `GET /api/users`: Tüm kullanıcıları listele (sadece admin)
- `GET /api/users/:id`: Belirli bir kullanıcıyı getir
- `PUT /api/users/:id`: Kullanıcı bilgilerini güncelle
- `DELETE /api/users/:id`: Kullanıcıyı sil (sadece admin)

### Ürün Endpoint'leri

- `POST /api/products`: Yeni ürün ekle
- `GET /api/products`: Tüm ürünleri listele
- `GET /api/products/:id`: Belirli bir ürünü getir
- `PUT /api/products/:id`: Ürün bilgilerini güncelle
- `DELETE /api/products/:id`: Ürünü sil

## Lisans

Bu proje MIT Lisansı ile lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın. 