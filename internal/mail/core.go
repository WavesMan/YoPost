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
	"os"
	"path/filepath"
	"strings"

	"github.com/YoPost/internal/config"
	"github.com/google/uuid"
)

type Core interface {
	ValidateUser(email string) bool
	GetConfig() *config.Config
	StoreEmail(from string, to []string, data string) error
}

type coreImpl struct {
	cfg *config.Config
}

func NewCore(cfg *config.Config) (Core, error) {
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

func (c *coreImpl) StoreEmail(from string, to []string, data string) error {
	if from == "" {
		return fmt.Errorf("empty from address")
	}
	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}
	if data == "" {
		return fmt.Errorf("empty email content")
	}

	// 使用系统临时目录存储邮件
	tmpDir := os.TempDir()
	emailDir := filepath.Join(tmpDir, "yopost_emails")
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

	return nil
}
