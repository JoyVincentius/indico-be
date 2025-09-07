# BE Assignment – Go + Gin

## Deskripsi
Service kecil yang menyediakan:
1. **Pemesanan produk dengan stok terbatas** – menjamin tidak ada overselling meski ada ribuan request bersamaan.
2. **Job settlement** – memproses jutaan transaksi, meng‑aggregate per merchant per hari, menulis hasil ke tabel `settlements` dan CSV yang dapat di‑download.

## Teknologi
- **Go 1.20+**
- **Gin** (HTTP router)
- **MySQL** (Docker Compose)
- **Docker Compose** untuk dev environment
- **Channels & Worker Pool** untuk background job
- **Context cancellation** untuk `POST /jobs/:id/cancel`

## Cara Menjalankan
```bash
# 1. Build & jalankan DB
docker compose up --build

# 2. Seed data (produk & transaksi)
go run scripts/seed_data.go

# 3. Jalankan API
go run main.go
