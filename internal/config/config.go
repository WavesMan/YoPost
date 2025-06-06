package config

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Auth   AuthConfig
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
