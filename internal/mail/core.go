package mail

import "github.com/YoPost/internal/config"

type Core interface {
	ValidateUser(email string) bool
	GetConfig() *config.Config
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
