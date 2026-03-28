#!/bin/bash
# Script Installer Otomatis INDVIM by Nasa

echo "🔍 Menambahkan Repository INDVIM (hastagaming)..."

# Tambahkan repo ke list sources Termux
REPO_URL="https://hastagaming.github.io/indvim-repo"
echo "deb [trusted=yes] $REPO_URL stable main" > $PREFIX/etc/apt/sources.list.d/indvim.list

echo "🔄 Updating package list..."
pkg update -y

echo "📦 Installing INDVIM..."
pkg install indvim -y

echo "✅ Selesai! Ketik 'indvim' untuk mulai ngoding"
