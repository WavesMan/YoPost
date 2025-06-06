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
