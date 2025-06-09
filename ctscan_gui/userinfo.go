package main

import "os/user"

type UserInfo struct {
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
