package protocol

import (
	"net"
	"strings"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type SMTPServer struct {
	cfg      *config.Config
	mailCore mail.Core
}

func NewSMTPServer(cfg *config.Config, mailCore mail.Core) *SMTPServer {
	return &SMTPServer{
		cfg:      cfg,
		mailCore: mailCore,
	}
}

func (s *SMTPServer) HandleCommand(conn net.Conn, cmd string) error {
	switch {
	case strings.HasPrefix(cmd, "EHLO"):
		_, err := conn.Write([]byte("250-HELO\r\n"))
		return err
	default:
		_, err := conn.Write([]byte("500 Unknown command\r\n"))
		return err
	}
}

func (s *SMTPServer) Start() error {
	// TODO: 实现SMTP服务器启动逻辑
	return nil
}
