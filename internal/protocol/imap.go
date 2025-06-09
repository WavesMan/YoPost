// IMAPServer 实现了IMAP协议服务端功能
//
// 主要功能:
// - 监听指定端口接收IMAP客户端连接
// - 处理基础IMAP命令交互
// - 支持通过上下文控制服务启停
//
// 使用NewIMAPServer创建实例，通过Start方法启动服务
// 可通过GetListener获取当前监听器实例
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

func (s *IMAPServer) Start(ctx context.Context) error {
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.IMAP.Port))
	log.Printf("INFO: 初始化IMAP服务 - 监听端口:%d, 超时设置:%v", 
		s.cfg.IMAP.Port, s.cfg.IMAP.Timeout)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("IMAP监听失败: %w", err)
	}
	s.listener = ln

	// 更新配置中的实际端口(当使用0时)
	if s.cfg.IMAP.Port == 0 {
		s.cfg.IMAP.Port = ln.Addr().(*net.TCPAddr).Port
	}

	log.Printf("IMAP服务启动成功，监听在 :%d\n", s.cfg.IMAP.Port)
	log.Printf("INFO: IMAP服务已就绪，等待客户端连接")

	// 使用通道处理服务器关闭
	serverClosed := make(chan struct{})
	defer close(serverClosed)

	go func() {
		select {
		case <-ctx.Done():
			ln.Close()
		case <-serverClosed:
			// 正常退出
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil // 正常关闭
			}
			return fmt.Errorf("接受IMAP连接失败: %w", err)
		}

		select {
		case <-ctx.Done():
			conn.Close()
			return nil
		default:
			go s.handleConnection(conn)
		}
	}
}

func (s *IMAPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("INFO: 新IMAP客户端连接 - 客户端:%s", conn.RemoteAddr().String())

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

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		args := parts[1:]

		log.Printf("INFO: 处理IMAP命令 - 命令:%s 参数:%v", command, args)

		switch command {
		case "LOGOUT":
			conn.Write([]byte("* BYE Logging out\r\n"))
			break
		case "SELECT":
			if len(args) < 1 {
				conn.Write([]byte("* BAD Missing mailbox name\r\n"))
				continue
			}
			conn.Write([]byte(fmt.Sprintf("* FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n* OK [PERMANENTFLAGS ()] Read-only mailbox\r\n* %d EXISTS\r\n* 0 RECENT\r\n* OK [UIDVALIDITY 1] UIDs valid\r\n* OK [UIDNEXT 1] Next UID\r\n* OK [READ-ONLY] Select completed\r\n", len(args))))
		case "FETCH":
			if len(args) < 2 {
				conn.Write([]byte("* BAD Missing FETCH arguments\r\n"))
				continue
			}
			// 简单实现：返回示例邮件内容
			conn.Write([]byte("* 1 FETCH (FLAGS (\\Seen) RFC822 {310}\r\n"))
			conn.Write([]byte("From: test@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nThis is a test email.\r\n)\r\n"))
		case "SEARCH":
			conn.Write([]byte("* SEARCH 1\r\n"))
		default:
			conn.Write([]byte("* BAD Unknown command\r\n"))
		}

		// 重置超时时间
		conn.SetDeadline(time.Now().Add(5 * time.Minute))
	}
}
