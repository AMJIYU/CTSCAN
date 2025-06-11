package pkg

import (
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp 创建一个新的 App 应用结构体
func NewApp() *App {
	return &App{}
}

// startup 在应用启动时被调用。保存 context 以便我们调用运行时方法
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}
