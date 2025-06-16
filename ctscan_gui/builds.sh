# 构建macos 64位版本
wails build  -platform  darwin/amd64 

# mac上构建windows 64位版本
wails build  -platform  windows/amd64 -skipbindings

# 构建linux 64位版本
wails build  -platform  linux/amd64 -skipbindings