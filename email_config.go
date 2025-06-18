package main

import (
	"os"
	"strconv"
)

// 邮件配置结构
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
}

// 获取邮件配置
func getEmailConfig() EmailConfig {
	// 从环境变量获取配置，如果没有则使用默认值
	host := getEnv("SMTP_HOST", "smtp.qq.com")
	portStr := getEnv("SMTP_PORT", "587")
	user := getEnv("SMTP_USER", "")
	pass := getEnv("SMTP_PASS", "")

	// 解析端口号
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587
	}

	// 如果邮箱配置为空，使用测试模式
	if user == "" || pass == "" {
		return EmailConfig{
			SMTPHost: "test",
			SMTPPort: 0,
			SMTPUser: "test@example.com",
			SMTPPass: "test",
		}
	}

	return EmailConfig{
		SMTPHost: host,
		SMTPPort: port,
		SMTPUser: user,
		SMTPPass: pass,
	}
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 