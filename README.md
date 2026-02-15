Notification Microservice (Go + gRPC + Couchbase)
Bu proje, bir sosyal medya platformu gereksinimleri doğrultusunda geliştirilmiş Bildirim Mikroservisi'dir. Servis, gRPC protokolü üzerinden iletişim sağlamakta ve veri kalıcılığı katmanında Couchbase (NoSQL) veritabanını kullanmaktadır.

Fonksiyonel Özellikler
SendNotification: Yeni bildirimlerin oluşturulması ve asenkron/senkron kayıt süreçlerinin yönetilmesi.

GetNotifications: Kullanıcıya özel bildirimlerin sayfalama (pagination) desteği ile listelenmesi.

MarkAsRead: Belirli bir bildirimin durumunun okundu olarak güncellenmesi.

GetUnreadCount: Kullanıcının toplam okunmamış bildirim sayısının hesaplanması.

Clean Architecture: Bağımlılıkların yönetimi için katmanlı mimari yapısı (Domain, Ports, Adapters).

Dockerization: Uygulama ve veritabanı bileşenlerinin Docker Compose ile orkestre edilmesi.

Teknik Yığın
Programlama Dili: Go (Golang) 1.25

İletişim Protokolü: gRPC & Protocol Buffers

Veri Depolama: Couchbase Server 7.2

Konteynerizasyon: Docker & Docker Compose

Test Stratejisi: Go Testing & Testify (Mocking)

Proje Mimarisi
Uygulama, Hexagonal Architecture (Ports and Adapters) prensipleri doğrultusunda, iş mantığını dış bileşenlerden izole edecek şekilde tasarlanmıştır:

internal/core/domain: İş kurallarını ve varlıkları (Notification entity) içerir. Herhangi bir dış kütüphane veya framework bağımlılığı bulunmaz.

internal/core/ports: Veritabanı ve servis katmanları için arayüz (interface) tanımlamalarını içerir.

internal/repository: Couchbase implementasyonunu barındıran veri erişim katmanıdır (Couchbase Adapter).

internal/service: İş mantığının (Business Logic) uygulandığı katmandır.

cmd: Uygulamanın giriş noktalarını (Server ve Client) barındırır.
Kurulum ve Çalıştırma
Sistemin çalıştırılması için çalışma ortamında Docker ve Docker Compose paketlerinin yüklü olması gerekmektedir.

Uygulamanın Başlatılması
Bash
docker-compose up --build
Not: Veritabanı servisinin hazır hale gelmesi ve bucket yapılandırmasının tamamlanması, sistem kaynaklarına bağlı olarak 15-30 saniye sürebilmektedir.

İstemci (Client) Testi
Servis aktif hale geldikten sonra, gRPC metotlarını test etmek için hazırda bulunan istemci uygulamasını kullanabilirsiniz:

Bash
go run cmd/client/main.go
Birim Testler
İş mantığı testleri, repository katmanı mocklanarak kurgulanmıştır. Testleri çalıştırmak için aşağıdaki komut kullanılır:

Bash
go test ./... -v
Teknik Kararlar ve Notlar
UUID: Dağıtık sistem mimarilerinde benzersizliği korumak ve çakışmaları engellemek adına kimlik yönetiminde UUID v4 tercih edilmiştir.

WaitUntilReady: Dağıtık ortamlarda uygulama ve veritabanı arasındaki "race condition" durumunu yönetmek amacıyla, bağlantı aşamasında bekleme ve yeniden deneme stratejisi uygulanmıştır.

NoSQL Seçimi: Bildirim verisinin esnek şema gereksinimi ve yüksek yazma (write-heavy) yükü göz önünde bulundurularak, düşük gecikme süreli Couchbase tercih edilmiştir.