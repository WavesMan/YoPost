// Core 定义了邮件服务的核心接口，包含用户验证、配置获取和邮件存储功能
//
// ValidateUser 验证用户邮箱是否有效（当前实现始终返回true）
// GetConfig 返回当前服务的配置信息
// StoreEmail 将邮件内容存储到系统临时目录中，包含发件人、收件人和邮件内容
//
//	参数:
//	  from - 发件人邮箱地址
//	  to - 收件人邮箱地址列表
//	  data - 邮件正文内容
//	返回:
//	  错误信息（如果存储过程中发生错误）
package mail

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/google/uuid"
)

type Core interface {
	ValidateUser(email string) bool
	GetConfig() *config.Config
	StoreEmail(from string, to []string, data string) error
	GetEmails() ([]Email, error)
	GetEmail(id string) (*Email, error)
}

type Email struct {
	ID      string
	From    string
	To      []string
	Subject string
	Date    string
	Body    string
	Read    bool
}

type coreImpl struct {
	cfg *config.Config
}

const (
	dataDir      = "data"
	emailsSubDir = "emails"
)

func ensureDataDir() error {
	path := filepath.Join(dataDir, emailsSubDir)
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create mail directory: %w", err)
	}
	return nil
}

func NewCore(cfg *config.Config) (Core, error) {
	if err := ensureDataDir(); err != nil {
		return nil, err
	}
	return &coreImpl{
		cfg: cfg,
	}, nil
}

func (c *coreImpl) ValidateUser(email string) bool {
	return true
}

func (c *coreImpl) GetConfig() *config.Config {
	return c.cfg
}

func (c *coreImpl) GetEmails() ([]Email, error) {
	emailDir := filepath.Join(dataDir, emailsSubDir)
	files, err := os.ReadDir(emailDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read email directory: %w", err)
	}

	var emails []Email
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".eml" {
			id := strings.TrimSuffix(file.Name(), ".eml")
			fileInfo, err := os.Stat(filepath.Join(emailDir, file.Name()))
			if err != nil {
				continue
			}
			emails = append(emails, Email{
				ID:   id,
				From: "",
				To:   nil,
				Date: fileInfo.ModTime().Format(time.RFC822),
			})
		}
	}
	return emails, nil
}

func (c *coreImpl) GetEmail(id string) (*Email, error) {
	emailPath := filepath.Join(dataDir, emailsSubDir, id+".eml")
	content, err := os.ReadFile(emailPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read email file: %w", err)
	}

	// 简单解析邮件内容
	parts := strings.SplitN(string(content), "\n\n", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid email format")
	}

	headers := strings.Split(parts[0], "\n")
	var from, to string
	for _, h := range headers {
		if strings.HasPrefix(h, "From: ") {
			from = strings.TrimPrefix(h, "From: ")
		} else if strings.HasPrefix(h, "To: ") {
			to = strings.TrimPrefix(h, "To: ")
		}
	}

	return &Email{
		ID:   id,
		From: from,
		To:   strings.Split(to, ","),
		Body: parts[1],
		Read: false,
	}, nil
}

func (c *coreImpl) StoreEmail(from string, to []string, data string) error {
	log.Printf("INFO: 开始存储邮件 - 发件人:%s, 收件人数量:%d", from, len(to))
	if from == "" {
		return fmt.Errorf("empty from address")
	}
	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}
	if data == "" {
		return fmt.Errorf("empty email content")
	}

	// 使用项目data目录存储邮件
	emailDir := filepath.Join(dataDir, emailsSubDir)
	emailID := uuid.New().String()
	filename := filepath.Join(emailDir, emailID+".eml")

	if err := os.MkdirAll(emailDir, 0755); err != nil {
		return fmt.Errorf("failed to create email directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create email file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("From: %s\nTo: %s\n\n%s", from, strings.Join(to, ","), data))
	if err != nil {
		return fmt.Errorf("failed to write email content: %w", err)
	}

	log.Printf("INFO: 邮件存储成功 - 邮件ID:%s, 存储路径:%s", emailID, filename)
	return nil
}
