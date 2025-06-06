package app

import (
	"context"
	"sync"

	"github.com/YoPost/internal/api"
	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
)

type Server struct {
	cfg        *config.Config
	httpServer *api.Server
	smtpServer *protocol.SMTPServer
	imapServer *protocol.IMAPServer
	pop3Server *protocol.POP3Server
	mailCore   mail.Core
}

func New(cfg *config.Config) (*Server, error) {
	// 初始化邮件核心
	mailCore, err := mail.NewCore(cfg)
	if err != nil {
		return nil, err
	}

	// 初始化API服务器
	httpServer, err := api.NewServer(cfg, mailCore)
	if err != nil {
		return nil, err
	}

	// 初始化协议服务器
	smtpServer := protocol.NewSMTPServer(cfg, mailCore)
	imapServer := protocol.NewIMAPServer(cfg, mailCore)
	pop3Server := protocol.NewPOP3Server(cfg, mailCore)

	return &Server{
		cfg:        cfg,
		httpServer: httpServer,
		smtpServer: smtpServer,
		imapServer: imapServer,
		pop3Server: pop3Server,
		mailCore:   mailCore,
	}, nil
}

func (s *Server) Start() error {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := s.httpServer.Start(); err != nil {
			// 处理错误
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.smtpServer.Start(); err != nil {
			// 处理错误
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.imapServer.Start(); err != nil {
			// 处理错误
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.pop3Server.Start(); err != nil {
			// 处理错误
		}
	}()

	wg.Wait()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	// 实现优雅关闭逻辑
	return nil
}
