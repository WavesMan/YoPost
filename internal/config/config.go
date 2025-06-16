package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// MailServerConfig 邮件服务器配置结构
type MailServerConfig struct {
	Mailserver struct {
		Host string `yaml:"host"`
		Smtp struct {
			TlsPort   string `yaml:"tls_port"`
			NotlsPort string `yaml:"notls_port"`
		} `yaml:"smtp"`
	} `yaml:"mailserver"`
}

// TestConfig 测试配置结构
type TestConfig struct {
	Userinfo []struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Email    string `yaml:"email"`
	} `yaml:"userinfo"`
	Email struct {
		TestSubject string `yaml:"test_subject"`
		TestBody    string `yaml:"test_body"`
	} `yaml:"email"`
}

// TLSConfig TLS配置结构
type TLSConfig struct {
	Enabled bool              `yaml:"enabled"`
	Servers map[string]string `yaml:"servers"` // 格式: "host:port": status
}

// LoadMailServerConfig 加载邮件服务器配置
func LoadMailServerConfig() (*MailServerConfig, error) {
	data, err := ioutil.ReadFile(filepath.Join("internal", "config", "mailserver.yml"))
	if err != nil {
		log.Printf("ERROR: Failed to read mailserver config - %v", err)
		return nil, err
	}

	var cfg MailServerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("ERROR: Failed to unmarshal mailserver config - %v", err)
		return nil, err
	}

	return &cfg, nil
}

// LoadTestConfig 加载测试配置
func LoadTestConfig() (*TestConfig, error) {
	data, err := ioutil.ReadFile(filepath.Join("internal", "config", "testconfig.yml"))
	if err != nil {
		log.Printf("ERROR: Failed to read test config - %v", err)
		return nil, err
	}

	var cfg TestConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("ERROR: Failed to unmarshal test config - %v", err)
		return nil, err
	}

	return &cfg, nil
}

// LoadTLSConfig 加载TLS配置
func LoadTLSConfig() (*TLSConfig, error) {
	data, err := ioutil.ReadFile(filepath.Join("internal", "config", "services.yml"))
	if err != nil {
		log.Printf("ERROR: Failed to read TLS config - %v", err)
		return nil, err
	}

	var cfg TLSConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("ERROR: Failed to unmarshal TLS config - %v", err)
		return nil, err
	}

	return &cfg, nil
}
