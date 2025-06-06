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
