package protocol

import (
	"fmt"
	"net"
	"strconv"
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
			_, err := conn.Write([]byte("250 OK\r\n"))
			return err
		}
	case "RCPT":
		if strings.HasPrefix(cmd, "RCPT TO:") {
			_, err := conn.Write([]byte("250 OK\r\n"))
			return err
		}
	case "DATA":
		_, err := conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))
		return err
	case "QUIT":
		_, err := conn.Write([]byte("221 Bye\r\n"))
		return err
	default:
		_, err := conn.Write([]byte("500 Unknown command\r\n"))
		return err
	}
	return nil
}

func (s *SMTPServer) Start() error {
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.SMTP.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听失败: %w", err)
	}
	defer ln.Close()

	fmt.Printf("SMTP服务监听在 :%d\n", s.cfg.SMTP.Port)

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

	conn.Write([]byte("220 YoPost SMTP Service Ready\r\n"))

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		cmd := string(buf[:n])
		s.HandleCommand(conn, cmd)
	}
}
