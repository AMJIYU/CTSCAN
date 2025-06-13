package main

import (
	"embed"
	"log"

	"ctscan_gui/pkg"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建一个 App 实例
	app, err := pkg.NewApp()
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}

	// 使用配置创建应用
	err = wails.Run(&options.App{
		Title:     "CTScan 应急响应工具箱", // 设置窗口标题
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("运行应用失败: %v", err)
	}
}
