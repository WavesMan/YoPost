package protocol

import (
	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type POP3Server struct {
	cfg      *config.Config
	mailCore mail.Core
}

func NewPOP3Server(cfg *config.Config, mailCore mail.Core) *POP3Server {
	return &POP3Server{
		cfg:      cfg,
		mailCore: mailCore,
	}
}

func (s *POP3Server) Start() error {
	// TODO: 实现POP3服务器启动逻辑
	return nil
}
