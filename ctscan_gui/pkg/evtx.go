package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xrawsec/golang-evtx/evtx"
)

// EVTXEvent 表示解析后的EVTX事件
type EVTXEvent struct {
	Time        string         `json:"time"`
	EventID     int            `json:"event_id"`
	Provider    string         `json:"provider"`
	Level       string         `json:"level"`
	Channel     string         `json:"channel"`
	Computer    string         `json:"computer"`
	UserID      string         `json:"user_id"`
	Description string         `json:"description"`
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

// ParseEVTXFile 解析EVTX文件
func (a *App) ParseEVTXFile(filePath string) ([]EVTXEvent, error) {
	log.Printf("开始解析EVTX文件: %s", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("文件不存在: %s", filePath)
		return nil, fmt.Errorf("文件不存在: %s", filePath)
	}

	// 检查文件扩展名
	if !strings.HasSuffix(strings.ToLower(filePath), ".evtx") {
		return nil, fmt.Errorf("文件格式错误：必须是 .evtx 文件")
	}

	// 打开EVTX文件
	log.Printf("正在打开EVTX文件...")
	ef, err := evtx.Open(filePath)
	if err != nil {
		log.Printf("打开文件失败: %v", err)
		if strings.Contains(err.Error(), "Corrupted header") {
			return nil, fmt.Errorf("文件格式错误：不是有效的 EVTX 文件或文件已损坏")
		}
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer ef.Close()
	log.Printf("成功打开EVTX文件")

	events := make([]EVTXEvent, 0)
	eventCount := 0

	// 解析所有事件
	log.Printf("开始解析事件...")
	for event := range ef.Events() {
		eventCount++
		log.Printf("正在解析第 %d 个事件", eventCount)

		// 创建路径
		providerPath := evtx.Path("System/Provider/@Name")
		levelPath := evtx.Path("System/Level")
		computerPath := evtx.Path("System/Computer")
		userIDPath := evtx.Path("System/Security/@UserID")
		eventDataPath := evtx.Path("EventData")
		systemPath := evtx.Path("System")
		userDataPath := evtx.Path("UserData")
		versionPath := evtx.Path("System/Version")
		qualifiersPath := evtx.Path("System/Qualifiers")
		taskPath := evtx.Path("System/Task")
		opcodePath := evtx.Path("System/Opcode")
		keywordsPath := evtx.Path("System/Keywords")
		processIDPath := evtx.Path("System/Execution/@ProcessID")
		threadIDPath := evtx.Path("System/Execution/@ThreadID")

		// 获取事件ID
		eventID := int(event.EventID())
		log.Printf("事件ID: %d", eventID)

		// 安全地获取各个字段
		provider := "未知"
		if p, err := event.GetString(&providerPath); err == nil && p != "" {
			provider = p
		} else {
			// 尝试从系统信息中获取提供者
			if sys, err := event.GetMap(&systemPath); err == nil {
				if providerMap, ok := (*sys)["Provider"]; ok {
					if providerName, ok := providerMap.(map[string]interface{})["@Name"]; ok {
						if name, ok := providerName.(string); ok {
							provider = name
						}
					}
				}
			}
		}
		log.Printf("提供者: %s", provider)

		level := "未知"
		if l, err := event.GetInt(&levelPath); err == nil {
			level = getEventLevel(int(l))
		} else {
			// 尝试从系统信息中获取级别
			if sys, err := event.GetMap(&systemPath); err == nil {
				if levelVal, ok := (*sys)["Level"]; ok {
					if levelInt, ok := levelVal.(int); ok {
						level = getEventLevel(levelInt)
					}
				}
			}
		}
		log.Printf("级别: %s", level)

		// 获取通道
		channel := event.Channel()
		log.Printf("通道: %s", channel)

		// 安全地获取计算机名
		computer := "未知"
		if c, err := event.GetString(&computerPath); err == nil && c != "" {
			computer = c
		}
		log.Printf("计算机: %s", computer)

		// 安全地获取用户ID
		userID := "未知"
		if u, err := event.GetString(&userIDPath); err == nil && u != "" {
			userID = u
		}
		log.Printf("用户ID: %s", userID)

		// 获取系统信息
		systemInfo := make(map[string]any)
		if sys, err := event.GetMap(&systemPath); err == nil {
			for key, value := range *sys {
				systemInfo[key] = value
			}
		}

		// 获取事件数据
		eventData := make(map[string]any)
		if data, err := event.GetMap(&eventDataPath); err == nil {
			for key, value := range *data {
				eventData[key] = value
			}
		}

		// 获取用户数据
		userData := make(map[string]any)
		if data, err := event.GetMap(&userDataPath); err == nil {
			for key, value := range *data {
				userData[key] = value
			}
		}

		// 获取消息
		message := ""
		messagePath := evtx.Path("System/EventData/Data[@Name='Message']")
		if msg, err := event.GetString(&messagePath); err == nil && msg != "" {
			message = msg
		}

		// 获取其他字段
		version := 0
		if v, err := event.GetInt(&versionPath); err == nil {
			version = int(v)
		}

		qualifiers := 0
		if q, err := event.GetInt(&qualifiersPath); err == nil {
			qualifiers = int(q)
		}

		task := 0
		if t, err := event.GetInt(&taskPath); err == nil {
			task = int(t)
		}

		opcode := 0
		if o, err := event.GetInt(&opcodePath); err == nil {
			opcode = int(o)
		}

		keywords := ""
		if k, err := event.GetString(&keywordsPath); err == nil && k != "" {
			keywords = k
		}

		processID := 0
		if p, err := event.GetInt(&processIDPath); err == nil {
			processID = int(p)
		}

		threadID := 0
		if t, err := event.GetInt(&threadIDPath); err == nil {
			threadID = int(t)
		}

		// 转换事件数据
		evt := EVTXEvent{
			Time:          event.TimeCreated().Format("2006-01-02 15:04:05"),
			EventID:       eventID,
			Provider:      provider,
			Level:         level,
			Channel:       channel,
			Computer:      computer,
			UserID:        userID,
			Description:   getEventDescription(event),
			Data:          make(map[string]any),
			EventRecordID: int(event.EventRecordID()),
			Version:       version,
			Qualifiers:    qualifiers,
			Task:          task,
			Opcode:        opcode,
			Keywords:      keywords,
			ProcessID:     processID,
			ThreadID:      threadID,
			Message:       message,
			SystemInfo:    systemInfo,
			EventData:     eventData,
			UserData:      userData,
		}

		events = append(events, evt)
	}

	log.Printf("解析完成，共解析 %d 个事件", len(events))
	return events, nil
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
	levelStr := "未知"
	switch level {
	case 1:
		levelStr = "严重"
	case 2:
		levelStr = "错误"
	case 3:
		levelStr = "警告"
	case 4:
		levelStr = "信息"
	case 5:
		levelStr = "详细"
	}
	log.Printf("事件级别转换: %d -> %s", level, levelStr)
	return levelStr
}

// getEventDescription 获取事件描述
func getEventDescription(event *evtx.GoEvtxMap) string {
	// 获取基本信息
	providerPath := evtx.Path("System/Provider/@Name")
	provider := "未知"
	if p, err := event.GetString(&providerPath); err == nil && p != "" {
		provider = p
	}

	// 获取事件ID
	eventID := event.EventID()

	// 获取事件数据
	eventDataPath := evtx.Path("EventData")
	description := fmt.Sprintf("事件ID: %d, 提供者: %s", eventID, provider)

	// 添加登入类型标注
	switch eventID {
	case 4624, 4648: // 登入成功
		// 检查是否是RDP登入
		logonTypePath := evtx.Path("EventData/Data[@Name='LogonType']")
		if logonType, err := event.GetInt(&logonTypePath); err == nil {
			switch logonType {
			case 2: // 交互式登入
				description += "\n登入类型: 本地交互式登入"
			case 3: // 网络登入
				description += "\n登入类型: 网络登入"
			case 4: // 批处理登入
				description += "\n登入类型: 批处理登入"
			case 5: // 服务登入
				description += "\n登入类型: 服务登入"
			case 7: // 解锁
				description += "\n登入类型: 工作站解锁"
			case 8: // 网络明文
				description += "\n登入类型: 网络明文登入"
			case 9: // 新凭证
				description += "\n登入类型: 新凭证登入"
			case 10: // 远程交互
				description += "\n登入类型: 远程交互式登入 (RDP)"
			case 11: // 缓存交互
				description += "\n登入类型: 缓存交互式登入"
			default:
				description += fmt.Sprintf("\n登入类型: 未知 (%d)", logonType)
			}
		}
	case 4625, 4647: // 登入失败
		// 检查是否是RDP登入
		logonTypePath := evtx.Path("EventData/Data[@Name='LogonType']")
		if logonType, err := event.GetInt(&logonTypePath); err == nil {
			switch logonType {
			case 2: // 交互式登入
				description += "\n登入类型: 本地交互式登入"
			case 3: // 网络登入
				description += "\n登入类型: 网络登入"
			case 4: // 批处理登入
				description += "\n登入类型: 批处理登入"
			case 5: // 服务登入
				description += "\n登入类型: 服务登入"
			case 7: // 解锁
				description += "\n登入类型: 工作站解锁"
			case 8: // 网络明文
				description += "\n登入类型: 网络明文登入"
			case 9: // 新凭证
				description += "\n登入类型: 新凭证登入"
			case 10: // 远程交互
				description += "\n登入类型: 远程交互式登入 (RDP)"
			case 11: // 缓存交互
				description += "\n登入类型: 缓存交互式登入"
			default:
				description += fmt.Sprintf("\n登入类型: 未知 (%d)", logonType)
			}
		}
	}

	// 尝试获取更多详细信息
	if eventData, err := event.GetMap(&eventDataPath); err == nil {
		// 遍历事件数据，添加到描述中
		for key, value := range *eventData {
			// 跳过一些不重要的字段
			if key == "SubjectUserSid" || key == "SubjectUserName" || key == "SubjectDomainName" {
				continue
			}
			description += fmt.Sprintf("\n%s: %v", key, value)
		}
	}

	// 尝试获取消息
	messagePath := evtx.Path("System/EventData/Data[@Name='Message']")
	if message, err := event.GetString(&messagePath); err == nil && message != "" {
		description += fmt.Sprintf("\n消息: %s", message)
	}

	log.Printf("事件描述: %s", description)
	return description
}
