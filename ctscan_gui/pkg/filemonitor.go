package pkg

import (
	"os"
	"runtime"
	"time"
)

type FileInfo struct {
	Path        string    `json:"path"`
	Exists      bool      `json:"exists"`
	Size        int64     `json:"size"`
	Mode        string    `json:"mode"`
	ModTime     time.Time `json:"mod_time"`
	CreateTime  time.Time `json:"create_time"`
	AccessTime  time.Time `json:"access_time"`
	ChangeTime  time.Time `json:"change_time"`
	IsDir       bool      `json:"is_dir"`
	IsSymlink   bool      `json:"is_symlink"`
	Owner       string    `json:"owner"`
	Group       string    `json:"group"`
	Permissions string    `json:"permissions"`
	Description string    `json:"description"`
}

// Linux 敏感文件列表
var linuxSensitiveFiles = []string{
	"/etc/passwd",
	"/etc/shadow",
	"/etc/sudoers",
	"/etc/hosts",
	"/etc/resolv.conf",
	"/etc/ssh/sshd_config",
	"/etc/ssh/ssh_config",
	"/etc/profile",
	"/etc/bashrc",
	"/etc/hosts.allow",
	"/etc/hosts.deny",
	"/etc/crontab",
	"/var/log/auth.log",
	"/var/log/secure",
	"/var/log/messages",
	"/var/log/syslog",
	"/var/log/wtmp",
	"/var/log/btmp",
	"/var/log/lastlog",
	"/var/log/audit/audit.log",
	"/etc/pam.d/common-auth",
	"/etc/pam.d/common-password",
	"/etc/pam.d/common-session",
	"/etc/pam.d/common-account",
	"/etc/security/access.conf",
	"/etc/security/limits.conf",
	"/etc/security/pwquality.conf",
	"/etc/security/faillock.conf",
	"/etc/security/opasswd",
	"/etc/security/console.perms",
	"/etc/security/console.handlers",
	"/etc/security/console.apps",
	"/etc/security/console.perms.d",
	"/etc/security/console.handlers.d",
	"/etc/security/console.apps.d",
}

// macOS 敏感文件列表
var macOSSensitiveFiles = []string{
	"/etc/passwd",
	"/etc/shadow",
	"/etc/sudoers",
	"/etc/hosts",
	"/etc/resolv.conf",
	"/etc/ssh/sshd_config",
	"/etc/ssh/ssh_config",
	"/etc/profile",
	"/etc/bashrc",
	"/etc/hosts.allow",
	"/etc/hosts.deny",
	"/etc/crontab",
	"/var/log/system.log",
	"/var/log/secure.log",
	"/var/log/audit.log",
	"/var/log/asl.log",
	"/var/log/install.log",
	"/var/log/wifi.log",
	"/var/log/fsck_hfs.log",
	"/var/log/fsck_apfs.log",
	"/Library/Logs/DiagnosticReports",
	"/Library/Logs/System.log",
	"/Library/Logs/AppleSystemLog",
	"/Library/Logs/AppleProfileService",
	"/Library/Logs/AppleProfileService/ProfileService.log",
	"/Library/Logs/AppleProfileService/ProfileService.error.log",
}

// Windows 敏感文件列表
var windowsSensitiveFiles = []string{
	"C:\\Windows\\System32\\config\\SAM",
	"C:\\Windows\\System32\\config\\SYSTEM",
	"C:\\Windows\\System32\\config\\SECURITY",
	"C:\\Windows\\System32\\config\\SOFTWARE",
	"C:\\Windows\\System32\\config\\DEFAULT",
	"C:\\Windows\\System32\\drivers\\etc\\hosts",
	"C:\\Windows\\System32\\drivers\\etc\\networks",
	"C:\\Windows\\System32\\drivers\\etc\\protocol",
	"C:\\Windows\\System32\\drivers\\etc\\services",
	"C:\\Windows\\System32\\drivers\\etc\\lmhosts.sam",
	"C:\\Windows\\System32\\winevt\\Logs\\Security.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\System.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Application.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Setup.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\ForwardedEvents.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-WindowsUpdateClient%4Operational.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-Windows Defender%4Operational.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-PowerShell%4Operational.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-TaskScheduler%4Operational.evtx",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-RemoteDesktopServices-RdpCoreTS%4Operational.evtx",
}

// 文件描述映射
var fileDescriptions = map[string]string{
	// Linux 系统文件描述
	"/etc/passwd":                    "用户账户信息文件",
	"/etc/shadow":                    "用户密码哈希文件",
	"/etc/sudoers":                   "sudo权限配置文件",
	"/etc/hosts":                     "主机名解析文件",
	"/etc/resolv.conf":               "DNS解析配置文件",
	"/etc/ssh/sshd_config":           "SSH服务器配置文件",
	"/etc/ssh/ssh_config":            "SSH客户端配置文件",
	"/etc/profile":                   "系统全局环境变量配置",
	"/etc/bashrc":                    "Bash shell配置文件",
	"/etc/hosts.allow":               "允许访问的主机配置",
	"/etc/hosts.deny":                "拒绝访问的主机配置",
	"/etc/crontab":                   "系统定时任务配置",
	"/var/log/auth.log":              "认证日志文件",
	"/var/log/secure":                "安全相关日志",
	"/var/log/messages":              "系统消息日志",
	"/var/log/syslog":                "系统日志",
	"/var/log/wtmp":                  "用户登录记录",
	"/var/log/btmp":                  "失败登录记录",
	"/var/log/lastlog":               "最后登录记录",
	"/var/log/audit/audit.log":       "审计日志",
	"/etc/pam.d/common-auth":         "PAM认证配置",
	"/etc/pam.d/common-password":     "PAM密码策略配置",
	"/etc/pam.d/common-session":      "PAM会话配置",
	"/etc/pam.d/common-account":      "PAM账户配置",
	"/etc/security/access.conf":      "访问控制配置",
	"/etc/security/limits.conf":      "系统资源限制配置",
	"/etc/security/pwquality.conf":   "密码质量配置",
	"/etc/security/faillock.conf":    "登录失败锁定配置",
	"/etc/security/opasswd":          "旧密码历史记录",
	"/etc/security/console.perms":    "控制台权限配置",
	"/etc/security/console.handlers": "控制台处理程序配置",
	"/etc/security/console.apps":     "控制台应用程序配置",

	// macOS 系统文件描述
	"/var/log/system.log":               "系统日志",
	"/var/log/secure.log":               "安全日志",
	"/var/log/audit.log":                "审计日志",
	"/var/log/asl.log":                  "Apple系统日志",
	"/var/log/install.log":              "安装日志",
	"/var/log/wifi.log":                 "WiFi连接日志",
	"/var/log/fsck_hfs.log":             "HFS文件系统检查日志",
	"/var/log/fsck_apfs.log":            "APFS文件系统检查日志",
	"/Library/Logs/DiagnosticReports":   "诊断报告目录",
	"/Library/Logs/System.log":          "系统日志",
	"/Library/Logs/AppleSystemLog":      "Apple系统日志",
	"/Library/Logs/AppleProfileService": "配置文件服务日志",

	// Windows 系统文件描述
	"C:\\Windows\\System32\\config\\SAM":                                                                       "安全账户管理器数据库",
	"C:\\Windows\\System32\\config\\SYSTEM":                                                                    "系统配置数据库",
	"C:\\Windows\\System32\\config\\SECURITY":                                                                  "安全配置数据库",
	"C:\\Windows\\System32\\config\\SOFTWARE":                                                                  "软件配置数据库",
	"C:\\Windows\\System32\\config\\DEFAULT":                                                                   "默认用户配置数据库",
	"C:\\Windows\\System32\\drivers\\etc\\hosts":                                                               "主机名解析文件",
	"C:\\Windows\\System32\\drivers\\etc\\networks":                                                            "网络配置",
	"C:\\Windows\\System32\\drivers\\etc\\protocol":                                                            "协议定义",
	"C:\\Windows\\System32\\drivers\\etc\\services":                                                            "服务端口定义",
	"C:\\Windows\\System32\\drivers\\etc\\lmhosts.sam":                                                         "LM主机文件示例",
	"C:\\Windows\\System32\\winevt\\Logs\\Security.evtx":                                                       "安全事件日志",
	"C:\\Windows\\System32\\winevt\\Logs\\System.evtx":                                                         "系统事件日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Application.evtx":                                                    "应用程序事件日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Setup.evtx":                                                          "安装事件日志",
	"C:\\Windows\\System32\\winevt\\Logs\\ForwardedEvents.evtx":                                                "转发事件日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-WindowsUpdateClient%4Operational.evtx":             "Windows更新日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-Windows Defender%4Operational.evtx":                "Windows Defender日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-PowerShell%4Operational.evtx":                      "PowerShell操作日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-TaskScheduler%4Operational.evtx":                   "任务计划程序日志",
	"C:\\Windows\\System32\\winevt\\Logs\\Microsoft-Windows-RemoteDesktopServices-RdpCoreTS%4Operational.evtx": "远程桌面服务日志",
}

// 获取当前平台的敏感文件列表
func getSensitiveFilesForCurrentPlatform() []string {
	switch runtime.GOOS {
	case "linux":
		return linuxSensitiveFiles
	case "darwin":
		return macOSSensitiveFiles
	case "windows":
		return windowsSensitiveFiles
	default:
		return []string{}
	}
}

// 获取敏感文件信息
func (a *App) GetSensitiveFileInfo() []FileInfo {
	sensitiveFiles := getSensitiveFilesForCurrentPlatform()
	fileInfos := make([]FileInfo, 0, len(sensitiveFiles))

	for _, path := range sensitiveFiles {
		info, err := os.Stat(path)
		// 如果文件不存在，跳过
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Path:   path,
			Exists: true,
		}

		// 添加文件描述
		if desc, ok := fileDescriptions[path]; ok {
			fileInfo.Description = desc
		}

		// 基本文件信息
		fileInfo.Size = info.Size()
		fileInfo.Mode = info.Mode().String()
		fileInfo.ModTime = info.ModTime()
		fileInfo.IsDir = info.IsDir()
		fileInfo.IsSymlink = info.Mode()&os.ModeSymlink != 0
		fileInfo.Permissions = info.Mode().Perm().String()

		// 获取系统特定的文件信息
		setPlatformFileInfo(info, &fileInfo)

		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos
}

// 获取系统特定的文件信息
func setPlatformFileInfo(info os.FileInfo, fileInfo *FileInfo) {
	fileInfo.CreateTime = info.ModTime()
	fileInfo.AccessTime = info.ModTime()
	fileInfo.ChangeTime = info.ModTime()
	fileInfo.Owner = "SYSTEM"
	fileInfo.Group = "Users"
}
