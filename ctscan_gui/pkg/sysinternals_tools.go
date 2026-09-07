package pkg

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	stdruntime "runtime"
	"strings"
)

const sysinternalsVendor = "Microsoft Sysinternals"

//go:embed tool_assets/thirds/*
var packagedSysinternals embed.FS

type SysinternalsTool struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	FileName      string `json:"file_name"`
	Description   string `json:"description"`
	BestFor       string `json:"best_for"`
	LocalPath     string `json:"local_path"`
	Available     bool   `json:"available"`
	Packaged      bool   `json:"packaged"`
	Supported     bool   `json:"supported"`
	RequiresAdmin bool   `json:"requires_admin"`
	Vendor        string `json:"vendor"`
}

type sysinternalsToolDefinition struct {
	id            string
	name          string
	fileName      string
	description   string
	bestFor       string
	requiresAdmin bool
}

var sysinternalsTools = []sysinternalsToolDefinition{
	{
		id:            "autoruns",
		name:          "Autoruns",
		fileName:      "Autoruns64.exe",
		description:   "微软官方自启动项排查工具，覆盖登录项、服务、驱动、计划任务、浏览器插件等持久化位置。",
		bestFor:       "排查自启动项和持久化点",
		requiresAdmin: true,
	},
	{
		id:            "procexp",
		name:          "Process Explorer",
		fileName:      "procexp64.exe",
		description:   "微软官方高级进程查看器，可查看进程树、句柄、DLL、签名、网络与父子进程关系。",
		bestFor:       "查看进程、句柄和 DLL",
		requiresAdmin: true,
	},
	{
		id:            "procmon",
		name:          "Process Monitor",
		fileName:      "Procmon64.exe",
		description:   "微软官方实时监控工具，可捕获文件系统、注册表、进程、线程和网络相关活动。",
		bestFor:       "实时进程、注册表和文件行为监控",
		requiresAdmin: true,
	},
}

func (a *App) GetSysinternalsTools() []SysinternalsTool {
	tools := make([]SysinternalsTool, 0, len(sysinternalsTools))
	for _, definition := range sysinternalsTools {
		tools = append(tools, a.sysinternalsTool(definition))
	}
	return tools
}

func (a *App) InstallSysinternalsTool(toolID string) (SysinternalsTool, error) {
	definition, err := getSysinternalsToolDefinition(toolID)
	if err != nil {
		return SysinternalsTool{}, err
	}

	roots, err := sysinternalsInstallDirs()
	if err != nil {
		return SysinternalsTool{}, err
	}

	failures := make([]string, 0, len(roots))
	for _, root := range roots {
		if err := os.MkdirAll(root, 0755); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", root, err))
			continue
		}

		targetPath := filepath.Join(root, definition.fileName)
		if err := extractPackagedSysinternalsExecutable(definition.fileName, targetPath); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", targetPath, err))
			continue
		}

		tool := a.sysinternalsTool(definition)
		tool.LocalPath = targetPath
		tool.Available = true
		return tool, nil
	}

	return SysinternalsTool{}, fmt.Errorf("释放工具失败: %s", strings.Join(failures, "; "))
}

func (a *App) LaunchSysinternalsTool(toolID string) error {
	definition, err := getSysinternalsToolDefinition(toolID)
	if err != nil {
		return err
	}
	if stdruntime.GOOS != "windows" {
		return fmt.Errorf("%s 只能在 Windows 上启动", definition.name)
	}

	tool := a.sysinternalsTool(definition)
	if !tool.Available {
		tool, err = a.InstallSysinternalsTool(toolID)
		if err != nil {
			return err
		}
	}
	if tool.LocalPath == "" {
		return fmt.Errorf("未找到 %s", definition.fileName)
	}

	if definition.requiresAdmin && packagedSysinternalsExecutableExists(definition.fileName) && shouldMigrateSysinternalsToolPath(tool.LocalPath, definition.fileName) {
		if migratedTool, installErr := a.InstallSysinternalsTool(toolID); installErr == nil && migratedTool.LocalPath != "" {
			tool = migratedTool
		}
	}

	return launchWindowsTool(tool.LocalPath, definition.requiresAdmin)
}

func (a *App) OpenSysinternalsToolsFolder() error {
	roots, err := sysinternalsInstallDirs()
	if err != nil {
		return err
	}

	failures := make([]string, 0, len(roots))
	for _, root := range roots {
		if err := os.MkdirAll(root, 0755); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", root, err))
			continue
		}
		return openPath(root)
	}

	return fmt.Errorf("打开工具目录失败: %s", strings.Join(failures, "; "))
}

func (a *App) sysinternalsTool(definition sysinternalsToolDefinition) SysinternalsTool {
	localPath := findSysinternalsExecutable(definition.fileName)
	return SysinternalsTool{
		ID:            definition.id,
		Name:          definition.name,
		FileName:      definition.fileName,
		Description:   definition.description,
		BestFor:       definition.bestFor,
		LocalPath:     localPath,
		Available:     localPath != "",
		Packaged:      packagedSysinternalsExecutableExists(definition.fileName),
		Supported:     stdruntime.GOOS == "windows",
		RequiresAdmin: definition.requiresAdmin,
		Vendor:        sysinternalsVendor,
	}
}

func getSysinternalsToolDefinition(toolID string) (sysinternalsToolDefinition, error) {
	for _, tool := range sysinternalsTools {
		if tool.id == toolID {
			return tool, nil
		}
	}
	return sysinternalsToolDefinition{}, fmt.Errorf("未知工具: %s", toolID)
}

func packagedSysinternalsExecutableExists(fileName string) bool {
	assetPath := filepath.ToSlash(filepath.Join("tool_assets", "thirds", fileName))
	info, err := fs.Stat(packagedSysinternals, assetPath)
	return err == nil && !info.IsDir()
}

func extractPackagedSysinternalsExecutable(fileName, targetPath string) error {
	assetPath := filepath.ToSlash(filepath.Join("tool_assets", "thirds", fileName))
	src, err := packagedSysinternals.Open(assetPath)
	if err != nil {
		return fmt.Errorf("内置资源中未找到 %s: %w", fileName, err)
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("创建工具目录失败: %w", err)
	}

	tmpPath := targetPath + ".tmp"
	dst, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("释放工具文件失败: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("释放工具文件失败: %w", err)
	}
	if err := dst.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("关闭工具文件失败: %w", err)
	}

	if err := replaceFile(tmpPath, targetPath); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}

func replaceFile(tmpPath, targetPath string) error {
	if err := os.Rename(tmpPath, targetPath); err == nil {
		return nil
	}

	backupPath := targetPath + ".bak"
	_ = os.Remove(backupPath)
	if fileExists(targetPath) {
		if err := os.Rename(targetPath, backupPath); err != nil {
			return fmt.Errorf("备份旧工具文件失败: %w", err)
		}
	}
	if err := os.Rename(tmpPath, targetPath); err != nil {
		if fileExists(backupPath) {
			_ = os.Rename(backupPath, targetPath)
		}
		return fmt.Errorf("保存工具文件失败: %w", err)
	}
	_ = os.Remove(backupPath)
	return nil
}

func findSysinternalsExecutable(fileName string) string {
	for _, dir := range sysinternalsSearchDirs() {
		path := filepath.Join(dir, fileName)
		if fileExists(path) {
			return path
		}
	}
	return ""
}

func sysinternalsSearchDirs() []string {
	dirs := make([]string, 0, 4)
	if installDirs, err := sysinternalsInstallDirs(); err == nil {
		dirs = append(dirs, installDirs...)
	}
	if executablePath, err := os.Executable(); err == nil {
		executableDir := filepath.Dir(executablePath)
		dirs = append(dirs, filepath.Join(executableDir, "thirds"))
	}
	if workingDir, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(workingDir, "thirds"))
	}
	return uniquePaths(dirs)
}

func sysinternalsCacheDir() (string, error) {
	roots, err := sysinternalsInstallDirs()
	if err != nil {
		return "", err
	}
	return roots[0], nil
}

func sysinternalsInstallDirs() ([]string, error) {
	dirs := make([]string, 0, 2)
	if stdruntime.GOOS == "windows" {
		if programData := os.Getenv("ProgramData"); programData != "" {
			dirs = append(dirs, filepath.Join(programData, "CTScan", "Sysinternals"))
		}
	}

	root, err := os.UserCacheDir()
	if err == nil && root != "" {
		dirs = append(dirs, filepath.Join(root, "CTScan", "Sysinternals"))
	}

	dirs = uniquePaths(dirs)
	if len(dirs) == 0 {
		if err != nil {
			return nil, fmt.Errorf("获取工具目录失败: %w", err)
		}
		return nil, fmt.Errorf("获取工具目录失败")
	}
	return dirs, nil
}

func shouldMigrateSysinternalsToolPath(currentPath, fileName string) bool {
	preferredDir, err := sysinternalsCacheDir()
	if err != nil || preferredDir == "" {
		return false
	}
	preferredPath := filepath.Join(preferredDir, fileName)
	return filepath.Clean(currentPath) != filepath.Clean(preferredPath)
}

func uniquePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if path == "" {
			continue
		}
		cleaned := filepath.Clean(path)
		key := cleaned
		if stdruntime.GOOS == "windows" {
			key = strings.ToLower(cleaned)
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, cleaned)
	}
	return result
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func openPath(path string) error {
	switch stdruntime.GOOS {
	case "windows":
		return exec.Command("explorer.exe", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	case "linux":
		return exec.Command("xdg-open", path).Start()
	default:
		return errors.New("当前系统不支持打开目录")
	}
}
