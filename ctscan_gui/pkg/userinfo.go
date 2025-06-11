package pkg

import (
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
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

type Win32_UserAccount struct {
	Name     string
	SID      string
	FullName string
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
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return nil
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT Name, SID, FullName FROM Win32_UserAccount")
	if err != nil {
		return nil
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	var users []SystemUser
	for i := 0; i < count; i++ {
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()

		nameVar, _ := oleutil.GetProperty(item, "Name")
		sidVar, _ := oleutil.GetProperty(item, "SID")
		fullNameVar, _ := oleutil.GetProperty(item, "FullName")

		homeDir := ""
		u, err := user.Lookup(nameVar.ToString())
		if err == nil {
			homeDir = u.HomeDir
		}

		users = append(users, SystemUser{
			Username: nameVar.ToString(),
			Uid:      sidVar.ToString(),
			Gid:      "0",
			HomeDir:  homeDir,
			Name:     fullNameVar.ToString(),
		})
		item.Release()
	}
	return users
}
