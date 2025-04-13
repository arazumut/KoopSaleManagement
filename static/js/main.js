/**
 * Kooperatif Ürün Satış ve Stok Takip Sistemi
 * Ana JavaScript Dosyası
 */

document.addEventListener('DOMContentLoaded', function() {
    // Toast mesajları için fonksiyon
    window.showToast = function(message, type = 'success', duration = 3000) {
        // Varsa önceki toast'u kaldır
        const existingToast = document.querySelector('.toast');
        if (existingToast) {
            existingToast.remove();
        }

        // Yeni toast oluştur
        const toast = document.createElement('div');
        toast.className = `toast toast-${type} animate-fade-in`;
        toast.innerHTML = `
            <div class="flex items-center">
                <span>${message}</span>
                <button class="ml-auto focus:outline-none" onclick="this.parentElement.parentElement.remove()">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                </button>
            </div>
        `;

        // Toast'u sayfaya ekle
        document.body.appendChild(toast);

        // Süre sonunda kaldır
        setTimeout(() => {
            toast.classList.add('opacity-0');
            setTimeout(() => {
                toast.remove();
            }, 300);
        }, duration);
    };

    // Form doğrulama
    window.validateForm = function(formId) {
        const form = document.getElementById(formId);
        let isValid = true;

        // Tüm zorunlu alanları kontrol et
        const requiredFields = form.querySelectorAll('[required]');
        requiredFields.forEach(field => {
            if (!field.value.trim()) {
                field.classList.add('border-red-500');
                isValid = false;
                
                // Hata mesajı göster
                const errorElement = field.nextElementSibling;
                if (errorElement && errorElement.classList.contains('error-message')) {
                    errorElement.textContent = 'Bu alan zorunludur';
                } else {
                    const error = document.createElement('div');
                    error.className = 'error-message text-red-500 text-sm mt-1';
                    error.textContent = 'Bu alan zorunludur';
                    field.insertAdjacentElement('afterend', error);
                }
            } else {
                field.classList.remove('border-red-500');
                
                // Hata mesajını temizle
                const errorElement = field.nextElementSibling;
                if (errorElement && errorElement.classList.contains('error-message')) {
                    errorElement.textContent = '';
                }
            }
        });

        return isValid;
    };
    
    // Sayısal alanlara sadece sayı girişi
    const numberInputs = document.querySelectorAll('input[type="number"]');
    numberInputs.forEach(input => {
        input.addEventListener('keypress', function(e) {
            if (!/[\d.,]/.test(e.key)) {
                e.preventDefault();
            }
        });
    });
    
    // Para formatı fonksiyonu
    window.formatCurrency = function(value) {
        return new Intl.NumberFormat('tr-TR', {
            style: 'currency', 
            currency: 'TRY',
            minimumFractionDigits: 2
        }).format(value);
    };
    
    // Tarih formatı fonksiyonu
    window.formatDate = function(date) {
        if (!date) return '';
        
        const d = new Date(date);
        return d.toLocaleDateString('tr-TR', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });
    };

    // Sayfa bazlı içerik yükleme
    window.loadContent = function(url, targetElementId) {
        const targetElement = document.getElementById(targetElementId);
        
        if (!targetElement) return;
        
        fetch(url)
            .then(response => response.text())
            .then(html => {
                targetElement.innerHTML = html;
            })
            .catch(error => {
                console.error('İçerik yüklenirken hata oluştu:', error);
                targetElement.innerHTML = '<p class="text-red-500 p-4">İçerik yüklenirken bir hata oluştu.</p>';
            });
    };

    // Tablo sütunlarına göre sıralama fonksiyonu
    window.sortTable = function(tableId, columnIndex) {
        const table = document.getElementById(tableId);
        const tbody = table.querySelector('tbody');
        const rows = Array.from(tbody.rows);
        
        // Sıralama yönünü belirle
        const sortDirection = table.getAttribute('data-sort-direction') === 'asc' ? 'desc' : 'asc';
        table.setAttribute('data-sort-direction', sortDirection);
        
        // Sütuna göre sırala
        rows.sort((a, b) => {
            const cellA = a.cells[columnIndex].textContent.trim();
            const cellB = b.cells[columnIndex].textContent.trim();
            
            // Sayısal değer kontrolü
            if (!isNaN(cellA) && !isNaN(cellB)) {
                return sortDirection === 'asc' 
                    ? parseFloat(cellA) - parseFloat(cellB)
                    : parseFloat(cellB) - parseFloat(cellA);
            }
            
            // Metin karşılaştırma
            return sortDirection === 'asc'
                ? cellA.localeCompare(cellB, 'tr')
                : cellB.localeCompare(cellA, 'tr');
        });
        
        // Tabloyu yeniden oluştur
        rows.forEach(row => tbody.appendChild(row));
    };

    // Barkod oluşturucu
    window.generateBarcode = function(value, targetElementId) {
        if (!window.JsBarcode) {
            console.error('JsBarcode kütüphanesi yüklenmemiş');
            return;
        }
        
        const element = document.getElementById(targetElementId);
        if (!element) return;
        
        JsBarcode(element, value, {
            format: "CODE128",
            lineColor: "#000",
            width: 2,
            height: 50,
            displayValue: true
        });
    };
    
    // Stok Uyarı Kontrolü
    window.checkStockAlerts = function() {
        fetch('/api/stock/critical')
            .then(response => response.json())
            .then(data => {
                if (data.length > 0) {
                    const alertIcon = document.querySelector('.stock-alert-icon');
                    if (alertIcon) {
                        alertIcon.classList.add('text-red-500');
                        alertIcon.classList.add('animate-pulse');
                    }
                }
            })
            .catch(error => console.error('Stok kontrolü yapılırken hata:', error));
    };
    
    // Sayfaya özgü JS'leri yükle
    const pageScript = document.querySelector('body').dataset.page;
    if (pageScript) {
        const script = document.createElement('script');
        script.src = `/static/js/pages/${pageScript}.js`;
        script.onerror = () => console.warn(`${pageScript}.js dosyası bulunamadı`);
        document.body.appendChild(script);
    }
    
    // Sayfa ilk yüklendiğinde stok uyarılarını kontrol et
    if (document.querySelector('.stock-alert-icon')) {
        checkStockAlerts();
    }
});

// Grafik çizimi fonksiyonu
window.drawChart = function(canvasId, type, labels, datasets, options = {}) {
    if (!window.Chart) {
        console.error('Chart.js kütüphanesi yüklenmemiş');
        return;
    }
    
    const ctx = document.getElementById(canvasId).getContext('2d');
    
    // Varsayılan seçenekler
    const defaultOptions = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'top',
            },
            tooltip: {
                mode: 'index',
                intersect: false,
            }
        }
    };
    
    // Grafik oluştur
    new Chart(ctx, {
        type: type,
        data: {
            labels: labels,
            datasets: datasets
        },
        options: { ...defaultOptions, ...options }
    });
}; 