#!/bin/bash

# ==========================================
# INDVIM GLOBAL INSTALLER (Anti-Alias Mode)
# Created by: Nasa (hastagaming) - Kediri
# ==========================================

echo "🚀 Memulai instalasi INDVIM secara global..."

# 1. Deteksi Lingkungan (Termux vs Linux)
if [ -d "$PREFIX/bin" ]; then
    BIN_DIR="$PREFIX/bin"
    CONF_FILE="$HOME/.bashrc"
    SUDO=""
    echo "📱 Terdeteksi: Lingkungan Termux"
else
    BIN_DIR="/usr/local/bin"
    CONF_FILE="$HOME/.bashrc"
    [ -f "$HOME/.zshrc" ] && CONF_FILE="$HOME/.zshrc"
    SUDO="sudo"
    echo "💻 Terdeteksi: Lingkungan Linux"
fi

# 2. Kompilasi Kode Go
echo "📦 Mengompilasi source code..."
go build -o indvim main.go
if [ $? -ne 0 ]; then
    echo "❌ Gagal mengompilasi! Pastikan Go sudah terinstal."
    exit 1
fi

# 3. Pindahkan ke System Path ($PATH)
echo "🚚 Memasang binary ke $BIN_DIR..."
$SUDO mv indvim $BIN_DIR/indvim
$SUDO chmod +x $BIN_DIR/indvim

# 4. PEMBERSIHAN TOTAL (Bagian yang kamu minta)
echo "🧹 Memeriksa dan menghapus alias yang mengganggu..."

# Hapus alias dari sesi aktif saat ini
unalias indvim 2>/dev/null

# Hapus baris alias dari .bashrc atau .zshrc secara otomatis
# Kita gunakan 'sed' untuk mencari baris yang mengandung 'alias indvim=' dan menghapusnya
if [ -f "$HOME/.bashrc" ]; then
    sed -i '/alias indvim=/d' "$HOME/.bashrc"
fi
if [ -f "$HOME/.zshrc" ]; then
    sed -i '/alias indvim=/d' "$HOME/.zshrc"
fi

# 5. Selesai
echo "--------------------------------------------------"
echo "✅ INDVIM BERHASIL TERPASANG!"
echo "Sekarang 'indvim' bisa dipanggil dari folder MANA SAJA."
echo "Catatan: Semua alias lama telah dihapus otomatis."
echo "Silakan ketik 'source $CONF_FILE' atau buka terminal baru."
echo "--------------------------------------------------"
