package pkg

import "strings"

// 从消息中提取用户名
func ExtractUsername(message string) string {
	if idx := strings.Index(message, "Account Name:"); idx != -1 {
		start := idx + len("Account Name:")
		end := strings.Index(message[start:], "\n")
		if end != -1 {
			return strings.TrimSpace(message[start : start+end])
		}
	}
	return ""
}

// 从消息中提取IP地址
func ExtractIPAddress(message string) string {
	if idx := strings.Index(message, "Source Network Address:"); idx != -1 {
		start := idx + len("Source Network Address:")
		end := strings.Index(message[start:], "\n")
		if end != -1 {
			return strings.TrimSpace(message[start : start+end])
		}
	}
	return ""
}
