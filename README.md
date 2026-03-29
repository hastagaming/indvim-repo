# 🚀 INDVIM (Indonesia Vim)
> Editor teks terminal super cepat, ringan, dan modern buatan **Nasa (hastagaming)**.

---

## ⌨️ Shortcuts & Kontrol
Gunakan kombinasi tombol ini untuk navigasi dan kontrol di dalam INDVIM:

| Tombol | Fungsi |
| :--- | :--- |
| `Panah (↑ ↓ ← →)` | Navigasi kursor ke seluruh teks |
| `Enter` | Membuat baris baru (Insert Mode) |
| `Backspace` | Menghapus karakter atau menggabungkan baris |
| `ESC` | Memunculkan **Peringatan Keluar** (Anti-Lupa) |
| **`Ctrl + B`** | **Save & Exit** (Simpan semua perubahan dan keluar) |

---

## ⚡ Power Snippets (Auto-Complete)
Ketik kode pemicu di bawah ini lalu tekan **SPASI** untuk memunculkan template otomatis:

| Kode | Hasil Template |
| :--- | :--- |
| `!kt` | Boilerplate Main Function **Kotlin** |
| `!java` | Struktur Class & Main Method **Java** |
| `!py` | Template Script **Python** (dengan `if __name__`) |
| `!html` | Struktur Dasar **HTML5** |

---

## 🎨 Fitur Utama
- **NvChad Vibe:** Tampilan UI dengan status bar berwarna dan nomor baris ala Neovim.
- **Syntax Highlighting:** Pewarnaan otomatis untuk kata kunci (`func`, `package`, `val`, dll).
- **Auto-Save:** Setiap ketikanmu langsung tersimpan aman ke file.
- **Smart Warning:** Jika kamu menekan `ESC`, INDVIM akan mengingatkanmu untuk keluar menggunakan `Ctrl + B` agar data tidak hilang.

---

## 📥 Cara Instal Instan (Termux)
Cukup copy dan paste perintah sakti ini di Termux kamu:

```bash
curl -s [https://raw.githubusercontent.com/hastagaming/indvim-repo/main/install.sh](https://raw.githubusercontent.com/hastagaming/indvim-repo/main/install.sh) | bash
```
---

## 🛠 Build from Source
If you want to compile INDVIM yourself:
1. Clone this repo.
2. Run `go build -o indvim main.go`.
3. Run `./indvim`.

---

## 🚀 Cara Update ke Versi Terbaru
Jika kamu melakukan **fork**, pastikan tetap sinkron dengan versi original (hastagaming) agar mendapatkan fitur terbaru:

1. Jalankan skrip update:
   ```bash
   ./update.sh
   ```
--

##catatan PENTING!!
Untuk menginstal INDVIM, cukup jalankan chmod +x install.sh && ./install.sh. Skrip ini akan mengonfigurasi INDVIM agar bisa berjalan secara global di sistem kamu!
