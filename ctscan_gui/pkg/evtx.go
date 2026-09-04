package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/0xrawsec/golang-evtx/evtx"
)

// EVTXEvent 表示解析后的EVTX事件
type EVTXEvent struct {
	Time        string         `json:"time"`
	TimeUTC     string         `json:"time_utc"`
	TimeLocal   string         `json:"time_local"`
	EventID     int            `json:"event_id"`
	Provider    string         `json:"provider"`
	Level       string         `json:"level"`
	Channel     string         `json:"channel"`
	Computer    string         `json:"computer"`
	UserID      string         `json:"user_id"`
	Description string         `json:"description"`
	EventType   string         `json:"event_type"`
	Data        map[string]any `json:"data"`
	// 新增字段
	EventRecordID int    `json:"event_record_id"`
	Version       int    `json:"version"`
	Qualifiers    int    `json:"qualifiers"`
	Task          int    `json:"task"`
	Opcode        int    `json:"opcode"`
	Keywords      string `json:"keywords"`
	ProcessID     int    `json:"process_id"`
	ThreadID      int    `json:"thread_id"`
	Message       string `json:"message"`
	// 系统信息
	SystemInfo map[string]any `json:"system_info"`
	// 事件数据
	EventData map[string]any `json:"event_data"`
	// 用户数据
	UserData map[string]any `json:"user_data"`
}

var (
	providerPath      = evtx.Path("/Event/System/Provider/Name")
	levelPath         = evtx.Path("/Event/System/Level")
	channelPath       = evtx.Path("/Event/System/Channel")
	computerPath      = evtx.Path("/Event/System/Computer")
	userIDPath        = evtx.Path("/Event/System/Security/UserID")
	eventDataPath     = evtx.Path("/Event/EventData")
	systemPath        = evtx.Path("/Event/System")
	userDataPath      = evtx.Path("/Event/UserData")
	versionPath       = evtx.Path("/Event/System/Version")
	qualifiersPath    = evtx.Path("/Event/System/Qualifiers")
	taskPath          = evtx.Path("/Event/System/Task")
	opcodePath        = evtx.Path("/Event/System/Opcode")
	keywordsPath      = evtx.Path("/Event/System/Keywords")
	processIDPath     = evtx.Path("/Event/System/Execution/ProcessID")
	threadIDPath      = evtx.Path("/Event/System/Execution/ThreadID")
	eventRecordIDPath = evtx.Path("/Event/System/EventRecordID")
	timeCreatedPath   = evtx.Path("/Event/System/TimeCreated/SystemTime")
)

// ParseEVTXFile 解析EVTX文件
func (a *App) ParseEVTXFile(filePath string) ([]EVTXEvent, error) {
	log.Printf("开始解析EVTX文件: %s", filePath)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件不存在: %s", filePath)
		}
		return nil, fmt.Errorf("无法访问文件 %s: %w", filePath, err)
	}

	if !strings.EqualFold(filepath.Ext(filePath), ".evtx") {
		return nil, fmt.Errorf("文件格式错误：必须是 .evtx 文件")
	}

	ef, err := evtx.Open(filePath)
	if err != nil {
		if strings.Contains(err.Error(), "Corrupted header") {
			return nil, fmt.Errorf("文件格式错误：不是有效的 EVTX 文件或文件已损坏")
		}
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer ef.Close()

	events := make([]EVTXEvent, 0, 4096)
	for event := range ef.Events() {
		events = append(events, parseEVTXEvent(event))
	}

	log.Printf("解析完成，共解析 %d 个事件", len(events))
	return events, nil
}

func parseEVTXEvent(event *evtx.GoEvtxMap) EVTXEvent {
	eventID := int(event.EventID())
	eventData := getMapValue(event, eventDataPath)
	provider := getStringValue(event, providerPath, "未知")
	level, err := event.GetInt(&levelPath)
	if err != nil {
		level = -1
	}

	timeCreated := ""
	timeCreatedUTC := ""
	timeCreatedLocal := ""
	if timestamp, err := event.GetTime(&timeCreatedPath); err == nil {
		timeCreated = timestamp.Local().Format("2006-01-02 15:04:05")
		timeCreatedUTC = timestamp.UTC().Format("2006-01-02 15:04:05")
		timeCreatedLocal = timestamp.Local().Format("2006-01-02 15:04:05")
	}

	message, _ := eventData["Message"].(string)
	return EVTXEvent{
		Time:          timeCreated,
		TimeUTC:       timeCreatedUTC,
		TimeLocal:     timeCreatedLocal,
		EventID:       eventID,
		Provider:      provider,
		Level:         getEventLevel(int(level)),
		Channel:       getStringValue(event, channelPath, "未知"),
		Computer:      getStringValue(event, computerPath, "未知"),
		UserID:        getStringValue(event, userIDPath, "未知"),
		Description:   getEventDescription(eventID, provider, eventData),
		EventType:     getEventType(eventID, eventData),
		Data:          make(map[string]any),
		EventRecordID: getIntValue(event, eventRecordIDPath),
		Version:       getIntValue(event, versionPath),
		Qualifiers:    getIntValue(event, qualifiersPath),
		Task:          getIntValue(event, taskPath),
		Opcode:        getIntValue(event, opcodePath),
		Keywords:      getStringValue(event, keywordsPath, ""),
		ProcessID:     getIntValue(event, processIDPath),
		ThreadID:      getIntValue(event, threadIDPath),
		Message:       message,
		SystemInfo:    getMapValue(event, systemPath),
		EventData:     eventData,
		UserData:      getMapValue(event, userDataPath),
	}
}

func getStringValue(event *evtx.GoEvtxMap, path evtx.GoEvtxPath, fallback string) string {
	value, err := event.GetString(&path)
	if err != nil || value == "" {
		return fallback
	}
	return value
}

func getIntValue(event *evtx.GoEvtxMap, path evtx.GoEvtxPath) int {
	value, err := event.GetInt(&path)
	if err != nil {
		return 0
	}
	return int(value)
}

func getMapValue(event *evtx.GoEvtxMap, path evtx.GoEvtxPath) map[string]any {
	result := make(map[string]any)
	value, err := event.Get(&path)
	if err != nil {
		return result
	}

	switch data := (*value).(type) {
	case evtx.GoEvtxMap:
		for key, item := range data {
			result[key] = item
		}
	case map[string]any:
		for key, item := range data {
			result[key] = item
		}
	}
	return result
}

// SaveEVTXFile 保存上传的EVTX文件
func (a *App) SaveEVTXFile(filePath string) (string, error) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "ctscan", "evtx")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 生成新的文件名
	newFileName := fmt.Sprintf("evtx_%d.evtx", time.Now().Unix())
	newFilePath := filepath.Join(tempDir, newFileName)

	// 复制文件
	srcFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开源文件失败: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(newFilePath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", fmt.Errorf("复制文件失败: %v", err)
	}

	return newFilePath, nil
}

// getEventLevel 获取事件级别
func getEventLevel(level int) string {
	switch level {
	case 0, 4:
		return "信息"
	case 1:
		return "严重"
	case 2:
		return "错误"
	case 3:
		return "警告"
	case 5:
		return "详细"
	default:
		return "未知"
	}
}

// getEventDescription 获取事件描述
func getEventDescription(eventID int, provider string, eventData map[string]any) string {
	description := fmt.Sprintf("事件ID: %d, 提供者: %s", eventID, provider)
	if eventID == 4624 || eventID == 4625 || eventID == 4648 {
		if logonType, ok := eventData["LogonType"]; ok {
			description += "\n登录类型: " + getLogonTypeDescription(fmt.Sprint(logonType))
		}
	}

	keys := make([]string, 0, len(eventData))
	for key := range eventData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if key == "SubjectUserSid" || key == "SubjectUserName" || key == "SubjectDomainName" || key == "Message" {
			continue
		}
		description += fmt.Sprintf("\n%s: %v", key, eventData[key])
	}

	if message, ok := eventData["Message"].(string); ok && message != "" {
		description += fmt.Sprintf("\n消息: %s", message)
	}
	return description
}

func getEventType(eventID int, eventData map[string]any) string {
	switch eventID {
	case 4624:
		return getSuccessfulLogonEventType(fmt.Sprint(eventData["LogonType"]))
	case 4625:
		return getFailedLogonEventType(fmt.Sprint(eventData["LogonType"]))
	case 4627:
		return "登录组成员信息"
	case 4634:
		return "账户注销"
	case 4647:
		return "用户主动注销"
	case 4648:
		return "显式凭据登录"
	case 4672:
		return "特权登录"
	case 4688:
		return "进程创建"
	case 4697:
		return "服务安装"
	case 4698:
		return "计划任务创建"
	case 4699:
		return "计划任务删除"
	case 4702:
		return "计划任务更新"
	case 4719:
		return "审计策略变更"
	case 4720:
		return "用户创建"
	case 4722:
		return "用户启用"
	case 4723:
		return "用户尝试修改密码"
	case 4724:
		return "用户密码重置"
	case 4725:
		return "用户禁用"
	case 4726:
		return "用户删除"
	case 4732:
		return "本地组成员添加"
	case 4733:
		return "本地组成员移除"
	case 4768:
		return "Kerberos TGT 请求"
	case 4769:
		return "Kerberos 服务票据请求"
	case 4771:
		return "Kerberos 预认证失败"
	case 4776:
		return "NTLM 凭据验证"
	case 4778:
		return "RDP 会话重新连接"
	case 4779:
		return "RDP 会话断开"
	case 4798:
		return "查询用户本地组成员"
	case 4799:
		return "查询启用安全组成员"
	case 5061:
		return "加密操作"
	case 1102:
		return "安全日志清除"
	default:
		return "普通安全事件"
	}
}

func getSuccessfulLogonEventType(logonType string) string {
	switch logonType {
	case "2":
		return "本地登录成功"
	case "3":
		return "网络登录成功"
	case "4":
		return "批处理登录成功"
	case "5":
		return "服务登录成功"
	case "7":
		return "工作站解锁成功"
	case "8":
		return "明文网络登录成功"
	case "9":
		return "新凭据登录成功"
	case "10":
		return "RDP 登录成功"
	case "11":
		return "缓存登录成功"
	default:
		return "登录成功"
	}
}

func getFailedLogonEventType(logonType string) string {
	if logonType == "" || logonType == "<nil>" {
		return "登录失败"
	}
	return "登录失败：" + getLogonTypeDescription(logonType)
}

func getLogonTypeDescription(logonType string) string {
	descriptions := map[string]string{
		"2":  "本地交互式登录",
		"3":  "网络登录",
		"4":  "批处理登录",
		"5":  "服务登录",
		"7":  "工作站解锁",
		"8":  "网络明文登录",
		"9":  "新凭据登录",
		"10": "远程交互式登录 (RDP)",
		"11": "缓存交互式登录",
	}
	if description, ok := descriptions[logonType]; ok {
		return description
	}
	return fmt.Sprintf("未知 (%s)", logonType)
}
