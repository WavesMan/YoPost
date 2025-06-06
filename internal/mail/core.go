// Package mail 提供邮件核心功能的接口和实现
//
// Core 接口定义了邮件服务的基本操作：
//   - ValidateUser: 验证用户邮箱有效性
//   - GetConfig: 获取当前配置
//   - StoreEmail: 存储邮件内容
//
// coreImpl 是 Core 接口的具体实现，包含配置信息
// NewCore 是创建 coreImpl 实例的构造函数
package mail

import "github.com/YoPost/internal/config"

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
	// TODO: 实现邮件存储逻辑
	return nil
}
