package protocol

import (
	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type SMTPServer struct {
	cfg      *config.Config
	mailCore *mail.Core
}

func NewSMTPServer(cfg *config.Config, mailCore *mail.Core) *SMTPServer {
	return &SMTPServer{
		cfg:      cfg,
		mailCore: mailCore,
	}
}

func (s *SMTPServer) Start() error {
	// TODO: 实现SMTP服务器启动逻辑
	return nil
}
