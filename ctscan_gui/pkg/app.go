package pkg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp 创建一个新的 App 应用结构体
func NewApp() (*App, error) {
	db, err := InitSqliteDB()
	if err != nil {
		return nil, err
	}
	return &App{db: db}, nil
}

// startup 在应用启动时被调用。保存 context 以便我们调用运行时方法
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectAndParseEVTXFile 弹窗选择EVTX文件并解析
func (a *App) SelectAndParseEVTXFile() ([]EVTXEvent, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择EVTX文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "Windows事件日志", Pattern: "*.evtx"},
		},
	})
	if err != nil {
		return nil, err
	}
	if filePath == "" {
		return nil, fmt.Errorf("未选择文件")
	}
	return a.ParseEVTXFile(filePath)
}
