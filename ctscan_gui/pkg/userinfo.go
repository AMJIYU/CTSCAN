package pkg

import (
	"os"
	"os/user"
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
