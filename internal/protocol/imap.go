package protocol

import (
	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type IMAPServer struct {
	cfg      *config.Config
	mailCore mail.Core
}

func NewIMAPServer(cfg *config.Config, mailCore mail.Core) *IMAPServer {
	return &IMAPServer{
		cfg:      cfg,
		mailCore: mailCore,
	}
}

func (s *IMAPServer) Start() error {
	// TODO: 实现IMAP服务器启动逻辑
	return nil
}
