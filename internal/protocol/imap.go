package protocol

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

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
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.IMAP.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("IMAP监听失败: %w", err)
	}
	defer ln.Close()

	fmt.Printf("IMAP服务监听在 :%d\n", s.cfg.IMAP.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("接受IMAP连接失败: %w", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *IMAPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 设置读写超时
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	// 发送欢迎消息
	if _, err := conn.Write([]byte("* OK YoPost IMAP Service Ready\r\n")); err != nil {
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

		// TODO: 实现完整IMAP命令处理
		if strings.EqualFold(cmd, "LOGOUT") {
			conn.Write([]byte("* BYE Logging out\r\n"))
			break
		}

		// 重置超时时间
		conn.SetDeadline(time.Now().Add(5 * time.Minute))
	}
}
