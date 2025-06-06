// Server 提供API服务核心功能
//
// cfg: 应用配置
// mailCore: 邮件服务核心组件
//
// NewServer 创建并初始化Server实例
// Start 启动API服务
package api

import (
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
	// TODO: 实现API服务器启动逻辑
	return nil
}
