```markdown
# Chat Moderasyon Servisi

Bu proje, Go dilinde yazılmış basit bir sohbet moderasyon servisidir. Bir MongoDB veritabanındaki sohbet günlüklerini izler, önceden tanımlanmış "küfürlü kelimeler" içeren mesajları tanımlar, bu olayları kaydeder ve moderasyon sistemi tarafından mesajları göz ardı edilen beyaz listeye alınmış oyunculara izin verir.

---

## Özellikler

* **Gerçek Zamanlı Sohbet İzleme:** Yeni sohbet mesajları için MongoDB'yi sürekli tarar.
* **Küfür Tespiti:** Yapılandırılabilir bir listeye göre uygunsuz dil içeren mesajları belirler.
* **Dinamik Küfür Listesi:** Bir dosyadaki değişiklikler algılandığında küfür listesini otomatik olarak yeniden yükler, servis yeniden başlatmaya gerek kalmaz.
* **Beyaz Listeye Alınmış Oyuncular:** Belirli oyuncuların moderasyon kontrollerini atlamasına izin verir.
* **CSV Kaydı:** Tespit edilen küfürleri, oyuncu adı, zaman damgası ve mesajla birlikte bir CSV dosyasına kaydeder.
* **Zarif Kapanma (Graceful Shutdown):** Sonlandırma sinyali alındığında tüm çalışan işlemlerin temiz bir şekilde kapanmasını sağlar.
* **Yapılandırılmış Loglama:** Verimli ve yapılandırılmış loglama için `zap` kullanır.

---

## Proje Yapısı

```
chat-moderation-service/
├── cmd/
│   └── main.go                  # Ana uygulama giriş noktası
├── config/
│   └── config.yaml              # Uygulama yapılandırma dosyası
├── internal/
│   ├── app/
│   │   ├── badwords.go          # Küfür listesini ve dosya izlemeyi yönetir
│   │   ├── filter.go            # Mesaj temizleme ve küfür tespit mantığını içerir
│   │   └── moderation.go        # MongoDB izleme ve moderasyon mantığını yönetir
│   ├── db/
│   │   └── mongo.go             # MongoDB bağlantı kurulumu
│   └── utils/
│       └── utils.go             # Yardımcı fonksiyonlar (yapılandırma yükleme)
└── badwords.txt                 # Küfürleri içeren dosya (her satırda bir tane)
```

---

## Başlarken

### Önkoşullar

* Go (1.16 veya üzeri sürüm önerilir)
* Çalışan bir MongoDB örneği (örn. `mongodb://localhost:27017`)

### Kurulum

1.  **Depoyu klonlayın:**
    ```bash
    git clone [https://github.com/kullanici-adiniz/chat-moderation-service.git](https://github.com/kullanici-adiniz/chat-moderation-service.git)
    cd chat-moderation-service
    ```
    (Yukarıdaki komutta `kullanici-adiniz` yerine kendi GitHub kullanıcı adınızı veya depo yolunuzu yazın.)
2.  **Bağımlılıkları indirin:**
    ```bash
    go mod tidy
    ```

### Yapılandırma

Uygulamanın yapılandırması `config/config.yaml` aracılığıyla yönetilir.

1.  **`config/config.yaml` dosyasını oluşturun:**
    ```yaml
    mongo_uri: "mongodb://localhost:27017" # MongoDB bağlantı URI'nız
    database_name: "foxlogger"             # MongoDB veritabanı adı
    collection_name: "logs"                # Sohbet günlüklerinin depolandığı MongoDB koleksiyon adı
    badwords_file: "badwords.txt"          # Küfür listesi dosyasının yolu
    allowed_players:                       # Beyaz listeye alınacak oyuncu isimleri listesi (büyük/küçük harf duyarlı)
      - ayd1ndemirci
      - DigerBeyazListeliOyuncu
    ```
2.  **`badwords.txt` dosyasını hazırlayın:**
    Proje kök dizininde bir `badwords.txt` dosyası oluşturun. Her küfürlü kelime veya ifade yeni bir satırda olmalıdır.
    ```
    küfür1
    başka kötü kelime
    salak
    ```

### Servisi Çalıştırma

```bash
go run cmd/main.go
```

Servis, belirtilen MongoDB koleksiyonunu izlemeye başlayacaktır. Tespit edilen küfürler proje kök dizinindeki `badwords_log.csv` dosyasına kaydedilecek ve konsola yazdırılacaktır.

---

## Nasıl Çalışır?

1.  **Başlatma (`cmd/main.go`):**
    * Yapılandırmayı `config/config.yaml` dosyasından yükler.
    * Yapılandırılmış loglama için `zap` logger'ı başlatır.
    * MongoDB'ye bir bağlantı kurar.
    * `badwords.txt` dosyasındaki değişiklikleri yüklemek ve izlemek için bir `BadWordManager` oluşturur.
    * MongoDB'ye bağlanan ve `BadWordManager`'ı kullanan bir `Moderator` başlatır.
    * Dosya izleme ve sohbet izleme için goroutine'leri başlatır.
    * Zarif kapanma için sinyal yakalamayı ayarlar.

2.  **Küfür Yönetimi (`internal/app/badwords.go`):**
    * `badwords.txt` dosyasındaki küfürleri belleğe okur.
    * `fsnotify` kullanarak `badwords.txt`'deki değişiklikleri algılar ve dosya değiştirildiğinde listeyi otomatik olarak yeniden yükler, böylece moderasyon sistemi yeniden başlatmaya gerek kalmadan güncel kalır.

3.  **Mesaj Filtreleme (`internal/app/filter.go`):**
    * `CleanMessage`: Sohbet mesajlarını küçük harfe dönüştürerek, alfasayısal olmayan karakterleri kaldırarak ve boşlukları kırparak ön işleme tabi tutar.
    * `ContainsBadWord`: Temizlenmiş bir mesajın `BadWordManager`'dan gelen herhangi bir küfür içerip içermediğini kontrol eder. **Not:** Mevcut uygulama, temizlendikten sonra `badwords.txt`'de görünen tek kelimeleri veya tam ifadeleri kontrol eder.

4.  **Moderasyon Mantığı (`internal/app/moderation.go`):**
    * Belirtilen MongoDB koleksiyonuna bağlanır.
    * Son kontrol edilen zaman damgasına göre sürekli olarak yeni `player_chat` olaylarını sorgular.
    * `allowed_players` listesindeki oyuncuların mesajlarını filtreler.
    * Beyaz listede olmayan mesajlar için küfür kontrolü yapar.
    * Tespit edilen küfürleri `badwords_log.csv` dosyasına kaydeder ve konsola bir uyarı yazdırır.
    * CSV dosyasını yönetir, dosya yeniyse başlıkları yazar ve günlükleri periyodik olarak diske kaydeder (flush).

---

## Gelecekteki İyileştirmeler (Değerlendirmeler)

* **Gelişmiş İfade Tespiti:** `ContainsBadWord` fonksiyonunu çok kelimeli küfürlü ifadeleri daha sağlam bir şekilde ele alacak şekilde geliştirmek, muhtemelen Aho-Corasick veya daha karmaşık regex desenleri kullanarak.
* **Hata Yönetimi ve Yeniden Denemeler:** Özellikle üretim ortamlarında MongoDB bağlantıları ve işlemleri için daha gelişmiş yeniden deneme mekanizmaları uygulamak.
* **Yapılandırılabilir Sorgulama Aralığı (Polling Interval):** `Moderator.StartMonitoring` içindeki izleme aralığının `config.yaml` aracılığıyla yapılandırılabilir olmasını sağlamak.
* **Metrikler ve İzleme:** Servis sağlığını, moderasyon sayılarını vb. izlemek için bir metrik sistemi (örn. Prometheus) ile entegrasyon.
* **Uyarı Sistemi:** Küfürler tespit edildiğinde uyarı gönderme işlevselliği eklemek (örn. Slack, e-posta yoluyla).

---
