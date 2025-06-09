package api

import (
	"log"
	"net/http"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type Server struct {
	cfg      *config.Config
	mailCore mail.Core
}

func NewServer(cfg *config.Config, mailCore mail.Core) (*Server, error) {
	return &Server{
		cfg:      cfg,
		mailCore: mailCore,
	}, nil
}

func (s *Server) Start() error {
    // 初始化路由器
    mux := http.NewServeMux()

    // 挂载邮件相关API路由
    mailHandler := mail.NewMailHandler(s.mailCore)
    mux.Handle("/api/mail/", http.StripPrefix("/api/mail", mailHandler))

    // 添加服务启动日志
    log.Printf("API服务启动中，监听地址: %s", s.cfg.Server.ListenAddr)
    log.Println("已注册路由:")
    log.Println("- POST /api/mail/smtp/send")

    // 启动HTTP服务器
    return http.ListenAndServe(s.cfg.Server.ListenAddr, mux)
}