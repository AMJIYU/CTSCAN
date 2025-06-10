package pkg

import (
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

type UserInfo struct {
	Username string `json:"username"`
	Uid      string `json:"uid"`
	Gid      string `json:"gid"`
	HomeDir  string `json:"home_dir"`
	Name     string `json:"name"`
}

type SystemUser struct {
	Username string `json:"username"`
	Uid      string `json:"uid"`
	Gid      string `json:"gid"`
	HomeDir  string `json:"home_dir"`
	Name     string `json:"name"`
}

func (a *App) GetUserInfo() UserInfo {
	u, err := user.Current()
	if err != nil {
		return UserInfo{}
	}
	return UserInfo{
		Username: u.Username,
		Uid:      u.Uid,
		Gid:      u.Gid,
		HomeDir:  u.HomeDir,
		Name:     u.Name,
	}
}

func (a *App) GetAllUsers() []SystemUser {
	switch runtime.GOOS {
	case "windows":
		return a.getWindowsUsers()
	case "linux", "darwin":
		return a.getUnixUsers()
	default:
		return nil
	}
}

func (a *App) getUnixUsers() []SystemUser {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	var users []SystemUser
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}
		users = append(users, SystemUser{
			Username: parts[0],
			Uid:      parts[2],
			Gid:      parts[3],
			HomeDir:  parts[5],
			Name:     parts[4],
		})
	}
	return users
}

func (a *App) getWindowsUsers() []SystemUser {
	// 使用 wmic 命令获取所有用户信息
	cmd := exec.Command("wmic", "useraccount", "get", "name,sid,fullname")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var users []SystemUser

	// 跳过标题行
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 解析用户信息
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// 获取用户主目录
		homeDir := ""
		u, err := user.Lookup(fields[0])
		if err == nil {
			homeDir = u.HomeDir
		}

		// 构建用户信息
		userInfo := SystemUser{
			Username: fields[0],
			Uid:      fields[1], // SID 作为 UID
			Gid:      "0",       // Windows 下使用默认 GID
			HomeDir:  homeDir,
			Name:     strings.Join(fields[2:], " "), // 全名可能包含空格
		}
		users = append(users, userInfo)
	}

	return users
}
