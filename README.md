# BU CASE STUDY 2NTech FİRMASI İÇİN HAZIRLANMIŞTIR.

# Notification Microservice (Go + gRPC + Couchbase)

Bu proje, yüksek performanslı ve ölçeklenebilir bir **Bildirim Mikroservisi**dir. Servis, **Hexagonal Architecture (Ports and Adapters)** prensiplerine sadık kalarak **Go** dilinde geliştirilmiş olup, veri kalıcılığı için **Couchbase** (NoSQL) kullanmaktadır.

Proje, modern mikroservis gereksinimlerini (Streaming, Rate Limiting, Observability) karşılayacak şekilde tasarlanmıştır.

## Özellikler ve Yetenekler

### 1. Temel Fonksiyonlar
- **SendNotification:** Yeni bildirim oluşturma, UUID atama ve veritabanına kaydetme.
- **GetNotifications:** Kullanıcı bildirimlerini sayfalama (pagination) ile listeleme (N1QL).
- **MarkAsRead:** Bildirimi okundu olarak işaretleme (Optimized `MutateIn` operation).
- **GetUnreadCount:** Okunmamış bildirim sayısını getirme (`COUNT` aggregation).

### 2. Bonus Özellikler
- **Server-Side Streaming (Live Updates):** `StreamNotifications` metodu ile kullanıcılara WebSocket benzeri canlı bildirim akışı sağlar. Go Channels ve Concurrency (Goroutines) kullanılarak implemente edilmiştir.
- **Rate Limiting (Throttling):** Token Bucket algoritması kullanılarak servisi aşırı yükten koruyan Middleware katmanı eklenmiştir.
- **Observability (Prometheus Metrics):** Sistem sağlığını ve anlık yükü izlemek için `/metrics` endpoint'i (Port: 9090) üzerinden Prometheus uyumlu veriler sunar.

### 3. Altyapı ve Mimari
- **Clean Architecture:** Domain, Service, Repository ve Transport katmanları tam izolasyonla ayrılmıştır.
- **Dockerized Environment:** Uygulama ve Veritabanı (Couchbase) `docker-compose` ile tek komutla ayağa kalkar.
- **Resiliency:** Veritabanı bağlantısı için `WaitUntilReady` ve `Healthcheck` mekanizmaları ile "Race Condition" önlenmiştir.

---

## Teknoloji Yığını

| Teknoloji | Amaç |
| :--- | :--- | :--- |
| **Go (Golang)** | Backend Dili |
| **gRPC & Protobuf** | İletişim Protokolü |
| **Couchbase** | NoSQL Veritabanı |
| **Docker & Compose** | Konteynerizasyon |
| **Testify** | Unit Test & Mocking |
| **Prometheus** | Metrik Toplama | Client Golang |

---

## Mimari Yapı (Clean Architecture)

Proje, bağımlılıkların dışarıdan içeriye doğru aktığı **Ports and Adapters** mimarisini kullanır.

notification-service/
├── cmd/
│   ├── client/          # Test İstemcisi (Streaming & CRUD testi)
│   └── server/          # Uygulama Giriş Noktası (Main, DI, Config)
├── internal/
│   ├── core/
│   │   ├── domain/      # Saf Entity'ler (Notification struct) - Bağımlılık YOK
│   │   └── ports/       # Interface'ler (Repository & Service Contracts)
│   ├── repository/      # Veritabanı Adaptörü (Couchbase Implementation)
│   ├── service/         # İş Mantığı (Business Logic & Use Cases)
│   └── middleware/      # Interceptorlar (Rate Limit, Metrics)
├── pkg/
│   └── pb/              # Generated gRPC Kodları (DTOs)
├── proto/               # Protobuf Kontrat Dosyası
└── docker-compose.yml   # Infrastructure Tanımı

---

## Kurulum ve Çalıştırma

Projeyi çalıştırmak için bilgisayarınızda **Docker** ve **Docker Compose** yüklü olmalıdır.

### 1. Servisi Başlatma
Tek bir komutla tüm altyapıyı (App + DB) ayağa kaldırın:

docker-compose up --build

Not: Couchbase ilk açılışta veritabanını ve indexleri hazırladığı için servisin tam yanıt vermesi 15-30 saniye sürebilir.

2. İstemci (Client) ile Test Etme
Case study kapsamında istenen Python scripti yerine, Go'nun Concurrency yeteneklerini sergileyen gelişmiş bir Client yazılmıştır. Bu istemci:

Arka planda (Goroutine) canlı yayın akışına (StreamNotifications) bağlanır.

Ana akışta bildirim gönderir (SendNotification).

Canlı yayına düşen bildirimi konsola basar.

Yeni bir terminalde çalıştırın:

go run cmd/client/main.go

3. Metrikleri İzleme (Prometheus)
Servis çalışırken tarayıcınızdan aşağıdaki adrese giderek anlık Goroutine sayısı, GC süresi ve İstek sayılarını görebilirsiniz:
http://localhost:9090/metrics

Testler
Business Logic, Repository katmanı Mocklanarak test edilmiştir. Veritabanına ihtiyaç duymadan iş kurallarını doğrular.

go test ./... -v

Docker Healthcheck:

Uygulamanın veritabanından önce başlayıp çökmesini (CrashLoopBackOff) engellemek için docker-compose tarafında Couchbase'in healthy durumu beklendi.
