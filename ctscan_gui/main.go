package main

import (
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"

	"ctscan_gui/pkg"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logFile := setupLogOutput()
	if logFile != nil {
		defer logFile.Close()
	}

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
		Windows:          windowsOptions(),
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("运行应用失败: %v", err)
	}
}

func setupLogOutput() *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logDir, err := os.UserCacheDir()
	if err != nil {
		logDir = os.TempDir()
	}
	logDir = filepath.Join(logDir, "CTScan", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil
	}

	logPath := filepath.Join(logDir, "ctscan.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	log.SetOutput(io.MultiWriter(os.Stderr, file))
	return file
}

func windowsOptions() *windows.Options {
	messages := windows.DefaultMessages()
	messages.InstallationRequired = "CTScan 需要 Microsoft Edge WebView2 Runtime。点击确定后请安装 WebView2 Runtime，然后重新打开 CTScan。"
	messages.UpdateRequired = "CTScan 需要更新 Microsoft Edge WebView2 Runtime。点击确定后请更新 WebView2 Runtime，然后重新打开 CTScan。"
	messages.MissingRequirements = "缺少运行环境"
	messages.Webview2NotInstalled = "未安装 Microsoft Edge WebView2 Runtime"
	messages.DownloadPage = "CTScan 需要 Microsoft Edge WebView2 Runtime。点击确定打开官方下载页面。最低版本要求："
	messages.ContactAdmin = "CTScan 需要 Microsoft Edge WebView2 Runtime 才能运行，请联系管理员安装。"
	messages.FailedToInstall = "WebView2 Runtime 安装失败，请手动安装后重新打开 CTScan。"
	messages.WebView2ProcessCrash = "WebView2 进程异常退出，请重启 CTScan；如果仍失败，请查看 CTScan 日志。"

	return &windows.Options{
		WebviewIsTransparent: false,
		WindowIsTranslucent:  false,
		DisableWindowIcon:    false,
		WebviewUserDataPath:  "",
		ZoomFactor:           1.0,
		Messages:             messages,
	}
}
