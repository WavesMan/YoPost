// Config 定义了应用程序的配置结构，包含服务器、数据库、认证、SMTP、IMAP和POP3的配置项
// Load 函数用于加载并返回配置实例，目前尚未实现具体逻辑
package config

import (
	"log"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Auth   AuthConfig
	SMTP   SMTPConfig
	IMAP   IMAPConfig
	POP3   POP3Config
}

type SMTPConfig struct {
	Port           int
	MaxMessageSize string
	Addr           string
	Domain         string
	MaxSize        int64
	TLSEnable      bool
	CertFile       string
	KeyFile        string
}

type IMAPConfig struct {
	Port    int
	Timeout string `yaml:"timeout"` // 新增超时配置项
}

type POP3Config struct {
	Port     int
	AuthType string `yaml:"auth_type"` // 新增认证类型配置项
}

type ServerConfig struct {
	Host       string `yaml:"host"`
	ListenAddr string `yaml:"listen_addr"` // 确保已有此字段
	Port       int    `yaml:"port"`
	Timeout    string `yaml:"timeout"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type AuthConfig struct {
	SecretKey string
	ExpiresIn int
}

func Load() (*Config, error) {
	log.Printf("INFO: 开始加载系统配置")
	return &Config{
		Server: ServerConfig{Host: "127.0.0.1", Port: 8080},
		SMTP: SMTPConfig{
			Port:           25,
			TLSEnable:      true,
			CertFile:       "cert.pem",
			KeyFile:        "key.pem",
			MaxMessageSize: "10MB",
		},
		IMAP: IMAPConfig{
			Port:    143,
			Timeout: "30s", // 添加默认超时配置
		},
		POP3: POP3Config{
			Port:     110,
			AuthType: "plain", // 添加默认认证类型
		},
	}, nil
}
