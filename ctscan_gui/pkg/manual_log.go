package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveManualLogFile 保存用户手动导出的页面日志。
func (a *App) SaveManualLogFile(defaultFileName string, content string) (string, error) {
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("日志内容为空")
	}

	fileName := sanitizeLogFileName(defaultFileName)
	if fileName == "" {
		fileName = "ctscan-log.html"
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".html"
	}

	defaultDir := defaultLogDirectory()
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		defaultDir = ""
	}
	filePath := ""
	if a.ctx != nil {
		selectedPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			Title:                "保存日志",
			DefaultDirectory:     defaultDir,
			DefaultFilename:      fileName,
			CanCreateDirectories: true,
			Filters: []runtime.FileFilter{
				{DisplayName: "HTML 报告 (*.html;*.htm)", Pattern: "*.html;*.htm"},
				{DisplayName: "Markdown 日志 (*.md)", Pattern: "*.md"},
				{DisplayName: "文本文件 (*.txt)", Pattern: "*.txt"},
				{DisplayName: "所有文件 (*.*)", Pattern: "*.*"},
			},
		})
		if err != nil {
			return "", fmt.Errorf("打开保存窗口失败: %w", err)
		}
		filePath = selectedPath
	} else {
		filePath = filepath.Join(defaultDir, fileName)
	}

	if strings.TrimSpace(filePath) == "" {
		return "", nil
	}
	if filepath.Ext(filePath) == "" {
		filePath += ".html"
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建保存目录失败: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("保存日志失败: %w", err)
	}
	return filePath, nil
}

func defaultLogDirectory() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		documentsDir := filepath.Join(homeDir, "Documents")
		if info, err := os.Stat(documentsDir); err == nil && info.IsDir() {
			return filepath.Join(documentsDir, "CTScan", "reports")
		}
		return filepath.Join(homeDir, "CTScan", "reports")
	}
	if documentsDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(documentsDir, "CTScan", "reports")
	}
	return filepath.Join(os.TempDir(), "CTScan", "reports")
}

func sanitizeLogFileName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '-'
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)
	name = strings.Trim(name, ". ")
	if len([]rune(name)) > 120 {
		runes := []rune(name)
		name = string(runes[:120])
	}
	return name
}
