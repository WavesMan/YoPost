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
	"bufio"
	"context"
	"crypto/tls" // 导入 crypto/tls 包以支持 TLS 配置
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
)

// SMTPConfig 包含SMTP服务器的配置参数
type SMTPConfig struct {
	Addr      string
	Domain    string
	MaxSize   int64
	TLSEnable bool
	CertFile  string
	KeyFile   string
	SMTP      struct {
		Port int
	}
}

type SMTPServer struct {
	config    *SMTPConfig
	mailCore  mail.Core
	tlsConfig *tls.Config
	conn      net.Conn
	listener  net.Listener
	reader    *bufio.Reader
	writer    *bufio.Writer
	state     sessionState
	cfg       *config.Config
}

// NewTestSMTPServer 创建用于测试的SMTP服务器实例
func NewTestSMTPServer(cfg *config.Config, mailCore mail.Core) (*SMTPServer, error) {
	return NewSMTPServer(cfg, mailCore)
}

type sessionState struct {
	sender     string
	recipients []string
	data       string
}

func (s *sessionState) Reset() {
	s.sender = ""
	s.recipients = s.recipients[:0]
	s.data = ""
}

func (s *SMTPServer) GetListener() net.Listener {
	return s.listener
}

func NewSMTPServer(cfg *config.Config, mailCore mail.Core) (*SMTPServer, error) {
	port := 25 // 默认非加密端口
	if cfg.SMTP.TLSEnable {
		port = 465 // 加密模式下使用465端口
	}

	server := &SMTPServer{
		config: &SMTPConfig{
			Addr:      fmt.Sprintf("%s:%d", cfg.Server.Host, port),
			Domain:    cfg.Server.Host,
			MaxSize:   cfg.SMTP.MaxSize,
			TLSEnable: cfg.SMTP.TLSEnable,
			CertFile:  cfg.SMTP.CertFile,
			KeyFile:   cfg.SMTP.KeyFile,
		},
		mailCore:  mailCore,
		tlsConfig: nil,
	}

	if cfg.SMTP.TLSEnable {
		cert, err := tls.LoadX509KeyPair(cfg.SMTP.CertFile, cfg.SMTP.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificate: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			PreferServerCipherSuites: true,
		}

		server.tlsConfig = tlsConfig
	}

	server.state = sessionState{
		sender:     "",
		recipients: make([]string, 0),
		data:       "",
	}

	return server, nil
}

func parsePortFromAddr(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 25 // 默认SMTP端口
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 25
	}
	return port
}

func (s *SMTPServer) Start(ctx context.Context) error {
	addr := s.config.Addr // 使用 s.config.Addr 替代 cfg.SMTP.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听失败: %w", err)
	}
	s.listener = ln
	defer ln.Close()

	log.Printf("SMTP服务监听在 %s\n", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受连接错误: %v", err)
			continue
		}
		go s.handleClient(conn)
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
			s.state.sender = strings.Trim(strings.TrimPrefix(cmd, "MAIL FROM:"), "<>")
			_, err := conn.Write([]byte("250 OK\r\n"))
			return err
		}
	case "RCPT":
		if strings.HasPrefix(cmd, "RCPT TO:") {
			to := strings.Trim(strings.TrimPrefix(cmd, "RCPT TO:"), "<>")
			s.state.recipients = append(s.state.recipients, to)
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

func (s *SMTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Reset state for new connection
	s.state.sender = ""
	s.state.recipients = s.state.recipients[:0]

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
				if err := s.mailCore.StoreEmail(s.state.sender, s.state.recipients, dataBuffer.String()); err != nil {
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

// handleClient 处理单个SMTP客户端连接
func (s *SMTPServer) handleClient(conn net.Conn) {
	defer conn.Close()

	s.conn = conn
	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)
	s.state.Reset()

	cmdHandlers := map[string]func(){
		"EHLO":     s.handleEHLO,
		"HELO":     s.handleHELO,
		"MAIL":     s.handleMAIL,
		"RCPT":     s.handleRCPT,
		"DATA":     s.handleDATA,
		"QUIT":     s.handleQUIT,
		"STARTTLS": s.handleSTARTTLS,
	}

	for {
		cmd, err := s.readCommand()
		if err != nil {
			return
		}

		handler, ok := cmdHandlers[cmd]
		if !ok {
			s.sendResponse("500 Unknown command\r\n")
			continue
		}

		handler()
	}
}

// readCommand 读取客户端发送的命令
func (s *SMTPServer) readCommand() (string, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// sendResponse 发送响应给客户端
func (s *SMTPServer) sendResponse(response string) {
	_, err := s.writer.WriteString(response)
	if err != nil {
		return
	}

	err = s.writer.Flush()
	if err != nil {
		return
	}
}

// handleEHLO 处理EHLO命令，初始化会话
func (s *SMTPServer) handleEHLO() {
	s.sendResponse("250-Hello\r\n")
	s.sendResponse("250-SIZE 10485760\r\n")
	if s.config.TLSEnable {
		s.sendResponse("250-STARTTLS\r\n")
	}
	s.sendResponse("250 OK\r\n")
}

// handleHELO 处理HELO命令，初始化会话
func (s *SMTPServer) handleHELO() {
	s.sendResponse("250 Hello\r\n")
}

// handleMAIL 处理MAIL命令
func (s *SMTPServer) handleMAIL() {
	cmd, err := s.readCommand()
	if err != nil {
		return
	}

	if !strings.HasPrefix(cmd, "MAIL FROM:") {
		s.sendResponse("501 Syntax error in parameters or arguments\r\n")
		return
	}

	s.state.sender = strings.Trim(strings.TrimPrefix(cmd, "MAIL FROM:"), "<>")
	s.sendResponse("250 OK\r\n")
}

// handleRCPT 处理RCPT命令
func (s *SMTPServer) handleRCPT() {
	cmd, err := s.readCommand()
	if err != nil {
		return
	}

	if !strings.HasPrefix(cmd, "RCPT TO:") {
		s.sendResponse("501 Syntax error in parameters or arguments\r\n")
		return
	}

	to := strings.Trim(strings.TrimPrefix(cmd, "RCPT TO:"), "<>")
	s.state.recipients = append(s.state.recipients, to)
	s.sendResponse("250 OK\r\n")
}

// handleDATA 处理DATA命令
func (s *SMTPServer) handleDATA() {
	s.sendResponse("354 End data with <CR><LF>.<CR><LF>\r\n")

	var data strings.Builder
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return
		}

		if line == ".\r\n" {
			break
		}

		// Remove leading dot if present (RFC 5321 section 4.5.2)
		if strings.HasPrefix(line, ".") {
			line = line[1:]
		}

		data.WriteString(line)
	}

	s.state.data = data.String()
	s.sendResponse("250 OK: Message accepted\r\n")

	if err := s.mailCore.StoreEmail(s.state.sender, s.state.recipients, s.state.data); err != nil {
		s.sendResponse("451 Requested action aborted: local error in processing\r\n")
		return
	}
}

// handleQUIT 处理QUIT命令
func (s *SMTPServer) handleQUIT() {
	s.sendResponse("221 Bye\r\n")
}

// handleSTARTTLS 处理STARTTLS命令，启动TLS加密连接
func (s *SMTPServer) handleSTARTTLS() {
	if !s.config.TLSEnable || s.tlsConfig == nil {
		s.sendResponse("421 TLS not available\r\n")
		return
	}

	s.sendResponse("220 Ready to start TLS\r\n")
	tlsConn := tls.Server(s.conn, s.tlsConfig)
	defer tlsConn.Close()

	// 重新初始化会话状态
	s.conn = tlsConn
	s.reader = bufio.NewReader(tlsConn)
	s.writer = bufio.NewWriter(tlsConn)
	s.state.Reset()

	// 增强TLS握手处理
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}
}

// GetState 返回当前会话状态，用于测试
func (s *SMTPServer) GetState() sessionState {
	return s.state
}

// ResetState 重置会话状态，用于测试
func (s *SMTPServer) ResetState() {
	s.state.Reset()
}

func (s *SMTPServer) listenAndServe() error {
	ln, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return err
	}
	s.listener = ln
	defer ln.Close()

	log.Printf("SMTP服务监听在 %s\n", s.config.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受连接错误: %v", err)
			continue
		}

		go s.handleClient(conn)
	}
}
