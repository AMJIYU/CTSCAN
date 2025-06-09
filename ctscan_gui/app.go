package main

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
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Stats 定义了概览统计的数据结构
type Stats struct {
	Total   int `json:"total"`
	Network int `json:"network"`
	Startup int `json:"startup"`
	Tasks   int `json:"tasks"`
	Process int `json:"process"`
}

// GetStats 返回初始的零值统计信息
// 这个方法被绑定了，所以可以从前端JS调用
func (a *App) GetStats() Stats {
	return Stats{
		Total:   10,
		Network: 0,
		Startup: 0,
		Tasks:   0,
		Process: 0,
	}
}
