// Config 定义了应用程序的配置结构，包含服务器、数据库、认证、SMTP、IMAP和POP3的配置项
// Load 函数用于加载并返回配置实例，目前尚未实现具体逻辑
package config

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
	Port int
}

type POP3Config struct {
	Port int
}

type ServerConfig struct {
    Host       string `yaml:"host"`
    ListenAddr string `yaml:"listen_addr"`  // 确保已有此字段
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
    SecretKey       string
    AllowedDomains  []string
    JWT             JWTConfig
    RateLimit       RateLimitConfig
    PasswordPolicy  PasswordPolicyConfig
}

// 新增JWT配置结构体
type JWTConfig struct {
    Expiration int    `yaml:"expiration"`
    Algorithm  string `yaml:"algorithm"`
}

// 新增速率限制配置
type RateLimitConfig struct {
    Enable   bool `yaml:"enable"`
    Requests int  `yaml:"requests"`
}

// 新增密码策略配置
type PasswordPolicyConfig struct {
    MinLength         int  `yaml:"min_length"`
    RequireMixedCase  bool `yaml:"require_mixed_case"`
    RequireNumbers    bool `yaml:"require_numbers"`
    RequireSymbols    bool `yaml:"require_symbols"`
}

func Load() (*Config, error) {
	// TODO: 实现配置加载逻辑
	return &Config{
		Server: ServerConfig{Host: "127.0.0.1", Port: 8080},
		SMTP: SMTPConfig{
			Port:           25,
			TLSEnable:      false,
			CertFile:       "cert.pem",
			KeyFile:        "key.pem",
			MaxMessageSize: "10MB",
		},
		IMAP: IMAPConfig{Port: 143},
		POP3: POP3Config{Port: 110},
	}, nil
}
