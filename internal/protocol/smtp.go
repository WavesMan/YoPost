// SMTPServer 实现了简单的SMTP协议服务器，用于接收和处理电子邮件
//
// 主要功能包括：
// - 监听指定端口接收SMTP连接
// - 处理标准SMTP命令（EHLO/HELO、MAIL FROM、RCPT TO、DATA、QUIT等）
// - 存储接收到的邮件到邮件核心系统
//
// 结构体包含配置信息、邮件核心处理模块和当前会话状态
package protocol

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

type SMTPServer struct {
	cfg         *config.Config
	mailCore    mail.Core
	currentFrom string
	currentTo   []string
	listener    net.Listener
}

func (s *SMTPServer) GetListener() net.Listener {
	return s.listener
}

func NewSMTPServer(cfg *config.Config, mailCore mail.Core) *SMTPServer {
	return &SMTPServer{
		cfg:      cfg,
		mailCore: mailCore,
	}
}

func (s *SMTPServer) HandleCommand(conn net.Conn, cmd string) error {
	cmd = strings.TrimSpace(cmd)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	verb := strings.ToUpper(parts[0])
	switch verb {
	case "EHLO", "HELO":
		_, err := conn.Write([]byte("250-yop.example.com\r\n250-PIPELINING\r\n250-8BITMIME\r\n250 SMTPUTF8\r\n"))
		return err
	case "MAIL":
		if strings.HasPrefix(cmd, "MAIL FROM:") {
			s.currentFrom = strings.Trim(strings.TrimPrefix(cmd, "MAIL FROM:"), "<>")
			_, err := conn.Write([]byte("250 OK\r\n"))
			return err
		}
	case "RCPT":
		if strings.HasPrefix(cmd, "RCPT TO:") {
			to := strings.Trim(strings.TrimPrefix(cmd, "RCPT TO:"), "<>")
			s.currentTo = append(s.currentTo, to)
			_, err := conn.Write([]byte("250 OK\r\n"))
			return err
		}
	case "DATA":
		// Handled in handleConnection
		return nil
	case "QUIT":
		_, err := conn.Write([]byte("221 Bye\r\n"))
		return err
	default:
		_, err := conn.Write([]byte("500 Unknown command\r\n"))
		return err
	}
	return nil
}

func (s *SMTPServer) Start(ctx context.Context) error {
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.SMTP.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听失败: %w", err)
	}
	s.listener = ln
	defer ln.Close()

	// 仅记录到日志，不输出到控制台
	log.Printf("SMTP服务监听在 :%d\n", s.cfg.SMTP.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("接受连接失败: %w", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *SMTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Reset state for new connection
	s.currentFrom = ""
	s.currentTo = nil

	// Send welcome message
	if _, err := conn.Write([]byte("220 YoPost SMTP Service Ready\r\n")); err != nil {
		return
	}

	inData := false
	var dataBuffer strings.Builder

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		text := string(buf[:n])
		text = strings.TrimRight(text, "\r\n")

		if inData {
			if text == "." {
				// End of DATA
				inData = false
				if _, err := conn.Write([]byte("250 OK: Message accepted\r\n")); err != nil {
					return
				}
				if err := s.mailCore.StoreEmail(s.currentFrom, s.currentTo, dataBuffer.String()); err != nil {
					conn.Write([]byte("451 Requested action aborted: local error in processing\r\n"))
					return
				}
				continue
			}

			// Remove leading dot if present (RFC 5321 section 4.5.2)
			if strings.HasPrefix(text, ".") {
				text = text[1:]
			}
			dataBuffer.WriteString(text + "\r\n")
			continue
		}

		cmd := strings.TrimSpace(text)
		if cmd == "" {
			continue
		}

		if cmd == "DATA" {
			if dataBuffer.Len() > 0 {
				dataBuffer.Reset()
			}
			inData = true
			if _, err := conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n")); err != nil {
				return
			}
			continue
		}

		if err := s.HandleCommand(conn, cmd); err != nil {
			return
		}
	}
}
