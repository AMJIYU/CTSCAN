#!/usr/bin/env bash
set -euo pipefail

if [ -d /opt/homebrew/opt/go/libexec/bin ]; then
  export PATH="/opt/homebrew/opt/go/libexec/bin:$PATH"
fi

# 构建 macOS 64 位版本
wails build -platform darwin/amd64

# macOS 上交叉构建 Windows 64 位版本。
# 使用静态链接，避免目标 Windows 缺少 MinGW/UCRT 运行库时双击无反应。
CGO_ENABLED=1 \
CC="${CC:-x86_64-w64-mingw32-gcc}" \
CGO_CFLAGS="${CGO_CFLAGS:--mcrtdll=msvcrt-os}" \
CGO_LDFLAGS="${CGO_LDFLAGS:--mcrtdll=msvcrt-os -static}" \
wails build \
  -platform windows/amd64 \
  -skipbindings \
  -webview2 error \
  -ldflags "-linkmode external -extldflags '-mcrtdll=msvcrt-os -static -static-libgcc -static-libstdc++'"

# 构建 Linux 64 位版本
wails build -platform linux/amd64 -skipbindings
