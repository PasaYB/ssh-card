# TUI Portfolio over SSH

Aplikasi portofolio interaktif berbasis Terminal User Interface (TUI) yang dikembangkan menggunakan Go dan pustaka [Wish](https://github.com/charmbracelet/wish) / [Bubble Tea](https://github.com/charmbracelet/bubbletea). Pengunjung dapat mengakses portofolio ini secara langsung melalui perintah SSH di terminal tanpa perlu melakukan registrasi, masukan password, atau instalasi dependen tambahan.

https://github.com/user-attachments/assets/1639ad2d-0f41-4ada-8246-d44f08d574fa

---

## Teknologi yang Digunakan

* **Go (Golang):** Bahasa pemrograman utama untuk logika aplikasi dan server SSH.
* **Wish (Charm):** Middleware SSH server khusus untuk aplikasi Go.
* **Bubble Tea (Charm):** Framework TUI berbasis arsitektur *The Elm Architecture*.
* **Lip Gloss (Charm):** Pustaka penataan gaya (*styling*) dan pewarnaan TUI.

---

## Konversi Gambar ke ASCII Art

Untuk menampilkan foto atau aset visual pada TUI, gambar berbasis piksel (`.jpeg`/`.png`) dikonversi menjadi format string teks ASCII menggunakan [ascii-image-converter](https://github.com/TheZoraiz/ascii-image-converter).

### Perintah Konversi
```bash
ascii-image-converter assets.jpeg -b -C -W 50 > assets/ascii.txt
```

## Panduan Deployment di Ubuntu

### 1. Pindahkan Port SSH Admin
Agar port `22` dapat digunakan oleh aplikasi Go TUI, ubah port OpenSSH server ke port lain (misalnya `2026`).

Buka file konfigurasi SSH Daemon:
```bash
sudo nano /etc/ssh/sshd_config
```
Ubah atau tambahkan baris berikut:
```bash
Port 2026
```
Pastikan port 2026 telah diizinkan pada Network Security Group (NSG) / Firewall sebelum menutup sesi SSH berjalan.

### 2. Kompilasi Kode Program
Masuk ke direktori proyek di server dan lakukan kompilasi:
```bash
cd /path/ke/projek-kamu
go build -o ssh-tui .
```

### 3. Membuat Systemd Service
Buat unit service agar aplikasi berjalan otomatis saat booting, mempertahankan izin port 22 (di bawah 1024) tanpa akses root penuh, serta memaksa output warna TrueColor.

Buat file service baru:
```bash
sudo nano /etc/systemd/system/tui.service
```
Isi dengan konfigurasi berikut (sesuaikan User dan WorkingDirectory):
```bash
[Unit]
Description=Go SSH TUI Portfolio Service
After=network.target

[Service]
Type=simple
User=azureuser
WorkingDirectory=/path/ke/projek-kamu
ExecStart=/path/ke/projek-kamu/ssh-tui

# Bind port < 1024 tanpa perlu akses root/sudo
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Memaksa TrueColor untuk sesi SSH
Environment="COLORTERM=truecolor"
Environment="TERM=xterm-256color"

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
### 4. Mengaktifkan dan Menjalankan Service
Muat ulang konfigurasi systemd, aktifkan auto-start, dan jalankan servicenya:
```bash
sudo systemctl daemon-reload
sudo systemctl enable tui
sudo systemctl start tui
```

## Cara Akses
Pengguna dapat mengakses portofolio ini dari terminal mana pun dengan menjalankan perintah:
```bash
ssh ssh.domainkamu.com
```
