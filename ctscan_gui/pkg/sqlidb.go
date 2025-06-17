package pkg

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/mattn/go-sqlite3" // 导入 SQLite 驱动程序
)

func getDesktopPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户主目录失败: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "Desktop"), nil
	case "darwin":
		return filepath.Join(homeDir, "Desktop"), nil
	case "linux":
		// Linux 系统可能有不同的桌面环境
		desktop := filepath.Join(homeDir, "Desktop")
		if _, err := os.Stat(desktop); os.IsNotExist(err) {
			// 尝试其他常见的桌面路径
			desktop = filepath.Join(homeDir, "桌面")
			if _, err := os.Stat(desktop); os.IsNotExist(err) {
				return "", fmt.Errorf("无法找到桌面目录")
			}
		}
		return desktop, nil
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

func InitSqliteDB() (*sql.DB, error) {
	desktopPath, err := getDesktopPath()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(desktopPath, "ctscan.db")

	// 检查数据库文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// 创建数据库文件
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("创建数据库文件失败: %v", err)
		}
		file.Close()
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %v", err)
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 创建表
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("创建表失败: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	// 创建用户信息表
	createUserInfoTable := `
	CREATE TABLE IF NOT EXISTS user_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		uid TEXT,
		gid TEXT,
		home_dir TEXT,
		name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建系统信息表
	createSystemInfoTable := `
	CREATE TABLE IF NOT EXISTS system_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname TEXT,
		os TEXT,
		arch TEXT,
		cpu_cores INTEGER,
		kernel_version TEXT,
		cpu_usage REAL,
		total_memory INTEGER,
		memory_usage REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建磁盘信息表
	createDiskInfoTable := `
	CREATE TABLE IF NOT EXISTS disk_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		system_info_id INTEGER,
		mount_point TEXT,
		total_size INTEGER,
		used_size INTEGER,
		free_size INTEGER,
		usage REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (system_info_id) REFERENCES system_info(id)
	);`

	// 创建定时任务表
	createCronTaskTable := `
	CREATE TABLE IF NOT EXISTS cron_task (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		line TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建敏感文件监控表
	createFileMonitorTable := `
	CREATE TABLE IF NOT EXISTS file_monitor (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT,
		file_exists BOOLEAN,
		size INTEGER,
		mode TEXT,
		mod_time DATETIME,
		create_time DATETIME,
		access_time DATETIME,
		change_time DATETIME,
		is_dir BOOLEAN,
		is_symlink BOOLEAN,
		owner TEXT,
		group_name TEXT,
		permissions TEXT,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建登录失败记录表
	createLoginFailedTable := `
	CREATE TABLE IF NOT EXISTS login_failed (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATETIME,
		event_id TEXT,
		event_type TEXT,
		source TEXT,
		username TEXT,
		ip_address TEXT,
		reason TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建登录成功记录表
	createLoginSuccessTable := `
	CREATE TABLE IF NOT EXISTS login_success (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATETIME,
		event_id TEXT,
		event_type TEXT,
		source TEXT,
		username TEXT,
		ip_address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建网络信息表
	createNetworkInfoTable := `
	CREATE TABLE IF NOT EXISTS network_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname TEXT,
		gateway TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建网络接口表
	createNetworkInterfaceTable := `
	CREATE TABLE IF NOT EXISTS network_interface (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		network_info_id INTEGER,
		name TEXT,
		ip TEXT,
		mac TEXT,
		bytes_sent INTEGER,
		bytes_recv INTEGER,
		packets_sent INTEGER,
		packets_recv INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (network_info_id) REFERENCES network_info(id)
	);`

	// 创建网络连接表
	createNetworkConnectionTable := `
	CREATE TABLE IF NOT EXISTS network_connection (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		proto TEXT,
		local_addr TEXT,
		remote_addr TEXT,
		status TEXT,
		pid INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建进程信息表
	createProcessInfoTable := `
	CREATE TABLE IF NOT EXISTS process_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pid INTEGER,
		name TEXT,
		ppid INTEGER,
		parent_name TEXT,
		create_time INTEGER,
		exe TEXT,
		file_ctime INTEGER,
		file_mtime INTEGER,
		md5 TEXT,
		signature TEXT,
		cpu_percent REAL,
		mem_percent REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建RDP登录表
	createRDPLoginTable := `
	CREATE TABLE IF NOT EXISTS rdp_login (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATETIME,
		username TEXT,
		ip TEXT,
		status TEXT,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建Shell历史记录表
	createShellHistoryTable := `
	CREATE TABLE IF NOT EXISTS shell_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATETIME,
		command TEXT,
		user TEXT,
		shell TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建启动项表
	createStartupItemTable := `
	CREATE TABLE IF NOT EXISTS startup_item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		path TEXT,
		type TEXT,
		enabled BOOLEAN,
		last_mod_time DATETIME,
		size INTEGER,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 执行创建表的SQL语句
	tables := []string{
		createUserInfoTable,
		createSystemInfoTable,
		createDiskInfoTable,
		createCronTaskTable,
		createFileMonitorTable,
		createLoginFailedTable,
		createLoginSuccessTable,
		createNetworkInfoTable,
		createNetworkInterfaceTable,
		createNetworkConnectionTable,
		createProcessInfoTable,
		createRDPLoginTable,
		createShellHistoryTable,
		createStartupItemTable,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("创建表失败: %v", err)
		}
	}

	return nil
}

// SaveUserInfo 保存用户信息到数据库
func (a *App) SaveUserInfo(userInfo UserInfo) error {
	query := `
	INSERT INTO user_info (username, uid, gid, home_dir, name)
	VALUES (?, ?, ?, ?, ?)`

	_, err := a.db.Exec(query,
		userInfo.Username,
		userInfo.Uid,
		userInfo.Gid,
		userInfo.HomeDir,
		userInfo.Name,
	)
	return err
}

// SaveSystemInfo 保存系统信息到数据库
func (a *App) SaveSystemInfo(sysInfo SystemInfo) error {
	// 开始事务
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入系统信息
	systemQuery := `
	INSERT INTO system_info (
		hostname, os, arch, cpu_cores, kernel_version,
		cpu_usage, total_memory, memory_usage
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(systemQuery,
		sysInfo.Hostname,
		sysInfo.OS,
		sysInfo.Arch,
		sysInfo.CPUCores,
		sysInfo.KernelVersion,
		sysInfo.CPUUsage,
		sysInfo.TotalMemory,
		sysInfo.MemoryUsage,
	)
	if err != nil {
		return err
	}

	// 获取插入的系统信息ID
	systemInfoID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// 插入磁盘信息
	diskQuery := `
	INSERT INTO disk_info (
		system_info_id, mount_point, total_size,
		used_size, free_size, usage
	) VALUES (?, ?, ?, ?, ?, ?)`

	for _, disk := range sysInfo.Disks {
		_, err = tx.Exec(diskQuery,
			systemInfoID,
			disk.MountPoint,
			disk.TotalSize,
			disk.UsedSize,
			disk.FreeSize,
			disk.Usage,
		)
		if err != nil {
			return err
		}
	}

	// 提交事务
	return tx.Commit()
}

// SaveCronTasks 保存定时任务到数据库
func (a *App) SaveCronTasks(tasks []CronTask) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO cron_task (line) VALUES (?)`
	for _, task := range tasks {
		_, err = tx.Exec(query, task.Line)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveFileMonitor 保存文件监控信息到数据库
func (a *App) SaveFileMonitor(files []FileInfo) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO file_monitor (
		path, file_exists, size, mode, mod_time,
		create_time, access_time, change_time, is_dir,
		is_symlink, owner, group_name, permissions, description
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, file := range files {
		_, err = tx.Exec(query,
			file.Path,
			file.Exists,
			file.Size,
			file.Mode,
			file.ModTime,
			file.CreateTime,
			file.AccessTime,
			file.ChangeTime,
			file.IsDir,
			file.IsSymlink,
			file.Owner,
			file.Group,
			file.Permissions,
			file.Description,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveLoginFailed 保存登录失败记录到数据库
func (a *App) SaveLoginFailed(records []LoginFailed) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO login_failed (
		time, event_id, event_type, source,
		username, ip_address, reason
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, record := range records {
		// 解析时间字符串
		timeValue, err := time.Parse("2006-01-02 15:04:05", record.Time)
		if err != nil {
			timeValue = time.Now() // 如果解析失败，使用当前时间
		}

		_, err = tx.Exec(query,
			timeValue,
			record.EventID,
			record.EventType,
			record.Source,
			record.Username,
			record.IPAddress,
			record.Reason,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveLoginSuccess 保存登录成功记录到数据库
func (a *App) SaveLoginSuccess(records []LoginSuccess) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO login_success (
		time, event_id, event_type, source,
		username, ip_address
	) VALUES (?, ?, ?, ?, ?, ?)`

	for _, record := range records {
		timeValue, err := time.Parse("2006-01-02 15:04:05", record.Time)
		if err != nil {
			timeValue = time.Now()
		}

		_, err = tx.Exec(query,
			timeValue,
			record.EventID,
			record.EventType,
			record.Source,
			record.Username,
			record.IPAddress,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveNetworkInfo 保存网络信息到数据库
func (a *App) SaveNetworkInfo(info NetworkInfo) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入网络基本信息
	networkQuery := `
	INSERT INTO network_info (hostname, gateway)
	VALUES (?, ?)`

	result, err := tx.Exec(networkQuery, info.Hostname, info.Gateway)
	if err != nil {
		return err
	}

	networkInfoID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// 插入网络接口信息
	interfaceQuery := `
	INSERT INTO network_interface (
		network_info_id, name, ip, mac,
		bytes_sent, bytes_recv, packets_sent, packets_recv
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	for i, iface := range info.Interfaces {
		var ip, mac string
		if i < len(info.IPs) {
			ip = info.IPs[i]
		}
		if i < len(info.MACs) {
			mac = info.MACs[i]
		}

		var stats InterfaceStats
		for _, stat := range info.InterfaceStats {
			if stat.Name == iface {
				stats = stat
				break
			}
		}

		_, err = tx.Exec(interfaceQuery,
			networkInfoID,
			iface,
			ip,
			mac,
			stats.BytesSent,
			stats.BytesRecv,
			stats.PacketsSent,
			stats.PacketsRecv,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveNetworkConnections 保存网络连接到数据库
func (a *App) SaveNetworkConnections(conns []NetworkConn) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO network_connection (
		proto, local_addr, remote_addr, status, pid
	) VALUES (?, ?, ?, ?, ?)`

	for _, conn := range conns {
		_, err = tx.Exec(query,
			conn.Proto,
			conn.LocalAddr,
			conn.RemoteAddr,
			conn.Status,
			conn.Pid,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveProcessInfo 保存进程信息到数据库
func (a *App) SaveProcessInfo(procs []ProcInfo) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO process_info (
		pid, name, ppid, parent_name, create_time,
		exe, file_ctime, file_mtime, md5, signature,
		cpu_percent, mem_percent
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, proc := range procs {
		_, err = tx.Exec(query,
			proc.PID,
			proc.Name,
			proc.PPID,
			proc.ParentName,
			proc.CreateTime,
			proc.Exe,
			proc.FileCtime,
			proc.FileMtime,
			proc.MD5,
			proc.Signature,
			proc.CPUPercent,
			proc.MemPercent,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveRDPLogin 保存RDP登录信息到数据库
func (a *App) SaveRDPLogin(logs []RDPLoginInfo) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO rdp_login (
		time, username, ip, status, description
	) VALUES (?, ?, ?, ?, ?)`

	for _, log := range logs {
		timeValue, err := time.Parse("2006-01-02 15:04:05", log.Time)
		if err != nil {
			timeValue = time.Now()
		}

		_, err = tx.Exec(query,
			timeValue,
			log.Username,
			log.IP,
			log.Status,
			log.Description,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveShellHistory 保存Shell历史记录到数据库
func (a *App) SaveShellHistory(records []ShellHistory) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO shell_history (
		time, command, user, shell
	) VALUES (?, ?, ?, ?)`

	for _, record := range records {
		timeValue, err := time.Parse("2006-01-02 15:04:05", record.Time)
		if err != nil {
			timeValue = time.Now()
		}

		_, err = tx.Exec(query,
			timeValue,
			record.Command,
			record.User,
			record.Shell,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SaveStartupItems 保存启动项到数据库
func (a *App) SaveStartupItems(items []StartupItem) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO startup_item (
		name, path, type, enabled, last_mod_time,
		size, description
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, item := range items {
		_, err = tx.Exec(query,
			item.Name,
			item.Path,
			item.Type,
			item.Enabled,
			item.LastModTime,
			item.Size,
			item.Description,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
