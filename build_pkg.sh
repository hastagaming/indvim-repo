#!/bin/bash
echo "😈 Memulai proses packaging INDVIM untuk Nasa..."

# 1. Compile file Go menjadi binary
go build -o indvim main.go

# 2. Buat struktur folder paket
mkdir -p indvim-pkg/data/data/com.termux/files/usr/bin
mkdir -p indvim-pkg/DEBIAN

# 3. Masukkan binary ke direktori usr/bin Termux
cp indvim indvim-pkg/data/data/com.termux/files/usr/bin/
chmod +x indvim-pkg/data/data/com.termux/files/usr/bin/indvim

# 4. Buat file Control (Metadata Package)
cat <<EOF > indvim-pkg/DEBIAN/control
Package: indvim
Version: 1.0.0
Architecture: all
Maintainer: Nasa
Description: INDVIM Text Editor - Super Cepat Tanpa Neovim
EOF

# 5. PERBAIKAN: Set permission yang diizinkan oleh dpkg-deb
chmod 755 indvim-pkg/DEBIAN
chmod 644 indvim-pkg/DEBIAN/control

# 6. Build folder menjadi file .deb
echo "Membungkus menjadi paket .deb..."
dpkg-deb --build indvim-pkg indvim.deb

# Bersihkan folder sementara agar rapi
rm -rf indvim-pkg

echo "✅ Paket indvim.deb berhasil dibuat tanpa error!"
