package protocol

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	listener net.Listener
}

func (s *IMAPServer) GetListener() net.Listener {
	return s.listener
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
	s.listener = ln
	defer ln.Close()

	log.Printf("IMAP服务监听在 :%d\n", s.cfg.IMAP.Port)

	// 使用context控制服务器生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// 如果是由于关闭导致的错误，不返回错误
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
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
