package mail

import "github.com/YoPost/internal/config"

type Core struct {
	cfg *config.Config
}

func NewCore(cfg *config.Config) (*Core, error) {
	return &Core{
		cfg: cfg,
	}, nil
}
