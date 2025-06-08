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
	Host string
	Port int
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
	// TODO: 实现配置加载逻辑
	return &Config{}, nil
}
