# Arama Motoru Servisi

Bu projede amaç, **farklı veri sağlayıcılardan (provider) alınan içerikleri** tek bir API altında toplayarak kullanıcıların anahtar kelime ile arama yapmasını sağlamak ve sonuçları skorlayarak sıralamaktır. Veriler JSON veya XML formatında olabilir; sistem bu farklı formatlardaki verileri normalize eder, puanlar ve kullanıcıya sunar. Ayrıca, projenin mimarisi yeni veri kaynaklarının kolayca eklenmesine ve mevcut kuralların değiştirilmesine imkân tanır.

![Frontend Screenshot](/static/ss.png)

## 🎯 Proje Kapsamı

**Case Çalışmasının Özet Gereksinimleri:**

- JSON ve XML formatında iki farklı sağlayıcıdan veri çekilmesi.
- Anahtar kelime araması, içerik türü filtresi (video/metin) ve skor bazlı sıralama.
- Standart bir puanlama algoritması ile popülerlik ve alaka sırasına göre sıralama.
- Kolay eklenebilir yeni provider mimarisi.
- Arama sonuçlarını basit bir web arayüzünde listeleme.
- Temiz kod, performans, hata yönetimi ve test edilebilirlik.

Bu gereksinimleri karşılamak için sistem şu başlıca bileşenlerden oluşur:

| Bileşen                          | Amaç                                                                                                |
| -------------------------------- | --------------------------------------------------------------------------------------------------- |
| **Backend (Go)**                 | İçeriklerin toplanması, normalize edilmesi, skorlama ve arama işlevlerinin API üzerinden sunulması. |
| **Veritabanı (PostgreSQL)**      | Kalıcı veri tutma: içerikler, metrikler, etiketler ve senkronizasyon logları.                       |
| **Cache (Redis)**                | Sık yapılan arama sorgularının ve sonuçların hızla yanıtlanabilmesi için ara bellek görevi.         |
| **Frontend (Basit Web Arayüzü)** | Kullanıcıların arama yapabilmesi ve sonuçları görüntüleyebilmesi.                                   |
| **Docker Compose**               | Tüm bileşenlerin (backend, frontend, Postgres, Redis) tek komutla çalıştırılması.                   |
| **Resilience**                   | **Rate Limiting** (Token Bucket) ve **Circuit Breaker** (Gobreaker) ile sistem kararlılığı.         |

## 🧱 Mimari Tasarım ve Kararlar

### Clean Architecture ve Katmanlar

Projede **Clean Architecture** yaklaşımı kullanıldı. Bu mimariyi tercih etmemin başlıca nedeni, **sistemi modüler hale getirerek bağımlılıkları kontrol altına almak** ve **kodun test edilebilirliğini ve genişletilebilirliğini artırmak**.

1. **Domain Katmanı**

   - Projenin en iç katmanında; `Content`, `ContentStats`, `Provider` ve `Tag` gibi temel veri modelleri yer alır.
   - `ScoringService` ile puanlama algoritması tek bir yerde tanımlanır. Bu servis, içerik türüne göre temel puanı hesaplar, güncellik ve etkileşim katsayılarını ekleyerek **final skor** üretir.
   - Domain katmanı sadece kendi iş kurallarıyla ilgilenir; veritabanı veya framework detaylarına bağlı değildir.

2. **Application (Use-Case) Katmanı**

   - `SearchContentsUseCase`: Anahtar kelime araması, filtreler ve sıralama kriterlerine göre arama yapar. Önce cache’e bakar; yoksa repository üzerinden veritabanından veriyi alır, skorlama yapar ve sonuçları sıralar.
   - `SyncProviderContentsUseCase`: Sağlayıcılardan verileri çekerek veritabanına senkronize eder. Provider’daki tüm içerikler periyodik olarak çekilir ve var olan kayıtlar güncellenir.
   - Bu katman, domain nesnelerini manipüle eder ve port arayüzlerini kullanarak dış dünya ile iletişime geçer.

3. **Infrastructure Katmanı**

   - Veritabanı erişimi için **sqlc** kullanıldı. Bu araç, SQL sorgularını Go kodu içerisinde derleme zamanında doğrulayarak tip güvenliğini ve performansı sağlar.
   - Sağlayıcılardan veri çekmek için `ProviderClient` arayüzü ve JSON/XML adaptörleri. Yeni bir format eklemek için bu arayüzü implemente etmek yeterlidir.
   - **Resilience**: `CircuitBreakerProviderClient` ile dış servis hatalarına karşı koruma sağlanır.
   - Redis cache adaptörü: Arama sonuçlarını anahtar bazlı saklamak için kullanılır.
   - Konfigürasyon: **Viper** ile dosya/env tabanlı konfigürasyon ve **DatabaseConfigProvider** ile veritabanı tabanlı dinamik skorlama kuralları yönetilir.

4. **Transport Katmanı**
   - **gRPC** sunucusu, düşük gecikme ve tip güvenliği sağlar.
   - **Rate Limiting**: gRPC interceptor ile API istekleri sınırlandırılır (Token Bucket algoritması).
   - **gRPC-Gateway** aracılığıyla aynı servisler HTTP/JSON olarak da kullanılabilir.
   - Basit bir web arayüzü, gRPC-Gateway üzerinden API’ye istek yapar.

### Neden Bu Teknolojileri Seçtik?

- **Go (Golang)**: Hafif, derlenmiş bir dil; eş zamanlı işlemleri kolay yönetir ve tek bir binary olarak dağıtım yapmaya olanak tanır. `net/http` paketi ve gRPC desteği güçlüdür.
- **gRPC + gRPC-Gateway**: gRPC hızlı ve güvenilir iken, gRPC-Gateway ile otomatik olarak REST benzeri JSON endpoint’ler elde edilir. Bu sayede tek bir servis tanımıyla hem performanslı gRPC hem de kolay kullanılır HTTP API sunuluyor.
- **PostgreSQL**: ACID uyumlu, gelişmiş veri tipleri ve full-text arama yetenekleri olan açık kaynak bir veritabanı. **Kalıcı veri tutarlılığı** için ideal.
- **sqlc**: ORM kullanmak yerine SQL sorgularını doğrudan yazıp tip güvenliği sağlamak için seçildi. Performans kaybı olmadan veritabanı işlemlerini yönetmek mümkün.
- **Redis**: Sık sorguların ve skoru hesaplanmış sonuçların çok hızlı döndürülmesini sağlamak için bellek içi cache kullanımı.
- **Docker Compose**: Production ortamında doğrudan kullanılmasa da bu case için tüm servisleri tek komutla ayağa kaldırmak amacıyla tercih edildi. Böylece kurulum süreci basitleşti.
- **React + Vite**: Hızlı ve modern bir frontend geliştirme deneyimi için React ile birlikte Vite build aracı kullanıldı. Tasarım minimal tutuldu.
- **GitHub Actions**: Sürekli entegrasyon (CI) süreçlerini otomatize etmek için kullanıldı. Her `push` ve `pull request` işleminde birim ve entegrasyon testleri otomatik olarak çalıştırılarak kodun kararlılığı sağlanır.

## 📊 Puanlama (Scoring) Algoritması

Case tanımında verilen puanlama formülü birebir uygulanmıştır:

\[ \text{Final Skor} = (\text{Temel Puan} \times \text{İçerik Türü Katsayısı}) + \text{Güncellik Puanı} + \text{Etkileşim Puanı} \]

- **Temel Puan**: Video için `views/1000 + likes/100`, metin için `reading_time + reactions/50`.
- **İçerik Türü Katsayısı**: Video için 1.5, metin için 1.0.
- **Güncellik Puanı**: İçeriğin yayın tarihine göre 1 hafta içinde +5, 1 ay içinde +3, 3 ay içinde +1 veya daha eski ise 0.
- **Etkileşim Puanı**: Video için `(likes/views) * 10` (views sıfırsa 0), metin için `(reactions/reading_time) * 5` (reading_time sıfırsa 0).

Bu bileşenler `ScoringService` içinde hesaplanır ve katsayılar veritabanındaki `scoring_rules` tablosundan dinamik olarak okunur. Bu sayede kod değişikliği yapmadan (deploy gerekmeden) puanlama algoritmasının ağırlıkları değiştirilebilir.

## 📦 Veri Yapısı

Sistem, verileri şu tablolarda saklar:

- `providers`: Sağlayıcı bilgileri (isim, format, URL, limit vb.).
- `contents`: İçerik metadata’sı (başlık, içerik türü, provider id, provider içerik id, yayın tarihi...).
- `content_stats`: İçeriklere ait ham metrikler (views, likes, reading_time, reactions, comments, duration_sec). Skor saklanmaz.
- `tags` & `content_tags`: Etiketlerin normalize edilmesi ve içeriklerle ilişkilendirilmesi.
- `content_raw_payloads`: (Opsiyonel) Orijinal JSON/XML verilerini saklama.
- `provider_sync_runs`: Sağlayıcı senkronizasyon işlemlerini ve loglarını takip etme.
- `scoring_rules`: Puanlama algoritması katsayılarını JSON formatında saklar.

Bu yapı, **kalıcı tutarlılık**, **normalize veri** ve **kolay genişletilebilirlik** sağlar. Ham veriler saklandığı için skorlama formülü değişse bile veriler yeniden işlenebilir.

## ⚙️ Kullanım Talimatları

Projeyi klonladıktan sonra hızlı bir şekilde çalıştırabilirsiniz:

1. **Depoyu Klonlayın:**
   ```bash
   git clone https://github.com/mehmetymw/search-aggregation-service.git
   cd search-aggregation-service
   ```
2. **Docker Compose ile Başlatın:**
   ```bash
   docker-compose up --build
   ```
   Bu komut PostgreSQL, Redis, backend ve frontend servislerini ayağa kaldıracaktır.
3. **Servisi Test Edin:**
   - API gRPC üzerinden `localhost:9090` portunda, HTTP/JSON üzerinden `http://localhost:8081` portunda çalışır.
   - Frontend arayüzü `http://localhost:5173` adresindedir.

### Örnek Arama İsteği

```
GET http://localhost:8081/api/v1/search?query=go%20programming&type=video&page=1&page_size=10
```

Yanıt:

```
{
  "items": [
    {
      "id": "42",
      "title": "Go Programming Tutorial",
      "content_type": "video",
      "score": 27.3,
      "published_at": "2024-03-15T10:00:00Z"
    },
    ...
  ],
  "page": 1,
  "page_size": 10,
  "total": 150
}
```

## 🔮 Geleceğe Yönelik İyileştirmeler

Proje case gereksinimlerini tamamen karşılıyor olsa da, gelecekte şu geliştirmelerle daha güçlü hale getirilebilir:

- **Tam Metin Arama:** PostgreSQL’in full-text search özelliklerini kullanarak daha iyi arama sonuçları.
- **Ölçeklenebilir Sync Mekanizması:** Şu anki in-memory ticker yerine, dağıtık sistemlerde sorunsuz çalışması için **CronJob** (Kubernetes) veya **Message Queue** (Kafka/RabbitMQ) tabanlı bir yapıya geçilebilir. Ancak bunun bir case study olması sebebiyle, daha basit bir yöntemle problemi çözmeye çalıştım, sync interval config.yaml dosyasında ayarlanabilir.
- **Arayüz İyileştirmeleri:** Kullanıcı deneyimini geliştirmek için daha interaktif filtreleme ve sıralama seçenekleri, grafiklerle zengin içerik.
- **Otomatik API Dokümantasyonu:** Swagger/OpenAPI entegrasyonu ile API’yi interaktif olarak belgelemek.

## 📝 Sonuç

Bu proje, verilen case çalışmasının tüm gereksinimlerini karşılamakla kalmayıp, **modüler ve ölçeklenebilir bir mimari** sunar. Kullanılan her teknoloji ve tasarım kararı, performans, esneklik ve bakım kolaylığı hedefleriyle uyumludur. Yeni provider’ların eklenmesi, puanlama formülünün değiştirilmesi veya arayüzün genişletilmesi gibi ihtiyaçlar basit dokunuşlarla gerçekleştirilebilir.

Bu README, projeyi sunarken kullanabileceğiniz temel noktaları özetler. Soru ve önerileriniz olursa memnuniyetle yanıtlarım.
