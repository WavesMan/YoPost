// POP3Server 实现了POP3协议服务器
//
// 主要功能包括：
// - 监听指定端口接收客户端连接
// - 处理基本的POP3命令交互
// - 提供邮件服务核心接口
//
// 使用NewPOP3Server创建实例，通过Start方法启动服务
package protocol

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type POP3Server struct {
	cfg      *config.Config
	mailCore mail.Core
	listener net.Listener
}

func (s *POP3Server) GetListener() net.Listener {
	return s.listener
}

func NewPOP3Server(cfg *config.Config, mailCore mail.Core) *POP3Server {
	return &POP3Server{
		cfg:      cfg,
		mailCore: mailCore,
	}
}

func (s *POP3Server) Start() error {
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.POP3.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("POP3监听失败: %w", err)
	}
	s.listener = ln
	defer ln.Close()

	log.Printf("POP3服务监听在 :%d\n", s.cfg.POP3.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("接受POP3连接失败: %w", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *POP3Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 设置读写超时
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	// 发送欢迎消息
	if _, err := conn.Write([]byte("+OK YoPost POP3 Service Ready\r\n")); err != nil {
		return
	}

	// 处理客户端命令
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		cmd := strings.TrimSpace(string(buf[:n]))
		if cmd == "" {
			continue
		}

		// TODO: 实现完整POP3命令处理
		if strings.EqualFold(cmd, "QUIT") {
			conn.Write([]byte("+OK Logging out\r\n"))
			break
		}

		// 重置超时时间
		conn.SetDeadline(time.Now().Add(5 * time.Minute))
	}
}
