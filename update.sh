#!/bin/bash
# launcher
chmod +x update.sh

# update.sh - Menarik update dari upstream (hastagaming)
echo "🔄 Mengambil update terbaru dari INDVIM Core..."
git remote add upstream https://github.com/hastagaming/indvim-repo.git 2>/dev/null
git fetch upstream
git checkout main
git merge upstream/main
go build -o indvim main.go
echo "✅ INDVIM berhasil diupdate ke versi terbaru!"
