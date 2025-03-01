
# Go Server Monitor

## Türkçe
### Proje Hakkında
Go Server Monitor, sunucu üzerindeki dosya yönetimi, süreç yönetimi, sistem bilgisi takibi, disk ve ağ verilerini izleme ile terminal yönetimi yapmanıza olanak tanıyan bir backend uygulamasıdır. Uygulama, Go programlama diliyle geliştirilmiş olup, çeşitli API endpoint'leri sunmaktadır.

### Mevcut Özellikler
- Dosya Yönetimi
- Süreç Yönetimi
- Sistem Bilgisi Takibi
- Disk ve Ağ Verilerini İzleme
- Terminal Yönetimi

#### Planlanan Ekstra Özellikler
- **Database Bağlantıları:** Veritabanı yönetimi ve veri saklama desteği (örn. PostgreSQL, MySQL)
- **Loglama:** Uygulama aktivitelerinin kaydedilmesi ve log analizi
- **Güvenlik ve Oturum Açma:** Kullanıcı oturum yönetimi ve kimlik doğrulama sistemleri
- **Endpointler için Token ile Doğrulama:** API isteklerinde JWT veya benzeri token doğrulama

### Kullanılan Teknolojiler
- Go 1.20
- Gorilla Mux (v1.8.1)
- Creack Pty (v1.1.24)
- Gorilla Websocket (v1.5.3)

### API Endpoint'leri

#### 1. Dosya Yönetimi
- **Dosya Listesi Getirme:** `GET /list-files?dir=/path`
- **Dosya Okuma:** `GET /read-file?file=/path/to/file`
- **Dosya İndirme:** `GET /download-file?file=/path/to/file`
- **Dosya Yazma/Güncelleme:** `POST /write-file` (Body: `{ "file": "/path", "content": "Yeni içerik" }`)
- **Dosya Silme:** `DELETE /delete-file?file=/path/to/file`

#### 2. Süreç Yönetimi
- **Süreçleri Listeleme:** `GET /processes`
- **Süreç Sonlandırma:** `POST /processes/kill?pid=1234`
- **Süreç Durdurma:** `POST /processes/stop?pid=1234`
- **Süreç Devam Ettirme:** `POST /processes/continue?pid=1234`
- **Süreç Yeniden Başlatma:** `POST /processes/restart?pid=1234`
- **Süreç Önceliğini Değiştirme:** `POST /processes/priority?pid=1234&priority=10`

#### 3. Sistem Bilgisi
- **Sistem Bilgisi Al:** `GET /system`

#### 4. Disk Bilgisi
- **Gerçek Zamanlı Disk Bilgisi:** `GET /disk` (WebSocket Bağlantısı)

#### 5. Ağ Bilgisi
- **Gerçek Zamanlı Ağ Bilgisi:** `GET /network` (WebSocket Bağlantısı)

#### 6. Terminal Yönetimi
- **Terminal Bağlantısı:** `GET /terminal` (WebSocket Bağlantısı)

### Kurulum ve Çalıştırma
```bash
# Projeyi klonlayın
git clone https://github.com/Gambitdutku/Go-server-monitor-

# Modülleri indirin
go mod download

# Uygulamayı çalıştırın
go run main.go

# Ya da compile ederek çalıştırın

GOOS=linux GOARCH=amd64 go build -o server_monitor

./server_monitor
```

### Katkıda Bulunanlar
- [Gambitdutku](https://github.com/Gambitdutku)

### Teşekkürler
Bu proje geliştirilirken aşağıdaki kütüphanelerden ve araçlardan yararlanılmıştır:
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Creack Pty](https://github.com/creack/pty)
- [Gorilla Websocket](https://github.com/gorilla/websocket)

---

## English
### About the Project
Go Server Monitor is a backend application that allows you to manage files, processes, monitor system information, disk and network data, and handle terminal management on a server. It is developed using the Go programming language and provides a set of API endpoints.

### Current and Planned Features
- File Management
- Process Management
- System Information Monitoring
- Disk and Network Data Monitoring
- Terminal Management

#### Planned Additional Features
- **Database Connections:** Support for database management and data storage (e.g., PostgreSQL, MySQL)
- **Logging:** Recording application activities and log analysis
- **Security and Authentication:** User session management and authentication systems
- **Token Validation for Endpoints:** Token validation for API requests using JWT or similar

### Technologies Used
- Go 1.20
- Gorilla Mux (v1.8.1)
- Creack Pty (v1.1.24)
- Gorilla Websocket (v1.5.3)

### API Endpoints

#### 1. File Management
- **Get File List:** `GET /list-files?dir=/path`
- **Read File:** `GET /read-file?file=/path/to/file`
- **Download File:** `GET /download-file?file=/path/to/file`
- **Write/Update File:** `POST /write-file` (Body: `{ "file": "/path", "content": "New content" }`)
- **Delete File:** `DELETE /delete-file?file=/path/to/file`

#### 2. Process Management
- **List Processes:** `GET /processes`
- **Kill Process:** `POST /processes/kill?pid=1234`
- **Stop Process:** `POST /processes/stop?pid=1234`
- **Continue Process:** `POST /processes/continue?pid=1234`
- **Restart Process:** `POST /processes/restart?pid=1234`
- **Change Process Priority:** `POST /processes/priority?pid=1234&priority=10`

#### 3. System Information
- **Get System Info:** `GET /system`

#### 4. Disk Information
- **Real-Time Disk Information:** `GET /disk` (WebSocket Connection)

#### 5. Network Information
- **Real-Time Network Information:** `GET /network` (WebSocket Connection)

#### 6. Terminal Management
- **Terminal Connection:** `GET /terminal` (WebSocket Connection)

### Installation and Run
```bash
# Clone the project
git clone https://github.com/Gambitdutku/Go-server-monitor-

# Download modules
go mod download

# Run the application
go run main.go

# Or run compiled application

GOOS=linux GOARCH=amd64 go build -o server_monitor

./server_monitor
```

### Contributors
- [Gambitdutku](https://github.com/Gambitdutku)

### Acknowledgements
Thanks to the following libraries and tools for their support in this project:
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Creack Pty](https://github.com/creack/pty)
- [Gorilla Websocket](https://github.com/gorilla/websocket)
