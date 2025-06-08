package protocol

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
)

type mockMailCore struct {
	storedEmails []struct {
		from       string
		recipients []string
		data       string
	}
}

func (m *mockMailCore) StoreEmail(from string, recipients []string, data string) error {
	m.storedEmails = append(m.storedEmails, struct {
		from       string
		recipients []string
		data       string
	}{
		from:       from, 
		recipients: recipients,
		data:       data,
	})
	return nil
}

func (m *mockMailCore) GetConfig() *config.Config {
	return &config.Config{}
}

func (m *mockMailCore) GetEmail(id string) (*mail.Email, error) {
	if len(m.storedEmails) > 0 {
		return &mail.Email{
			ID:      id,
			Content: m.storedEmails[0].data,
			From:    m.storedEmails[0].from,
			To:      m.storedEmails[0].recipients,
		}, nil
	}
	return nil, fmt.Errorf("email not found")
}

// 添加缺失的GetEmails方法实现
func (m *mockMailCore) GetEmails() ([]*mail.Email, error) {
	emails := make([]*mail.Email, len(m.storedEmails))
	for i, stored := range m.storedEmails {
		emails[i] = &mail.Email{
			ID:      fmt.Sprintf("email-%d", i),
			Content: stored.data,
			From:    stored.from,
			To:      stored.recipients,
		}
	}
	return emails, nil
}

func TestSMTPServer(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
		},
		SMTP: config.SMTPConfig{
			Port:           1025,
			Addr:           "127.0.0.1:1025",
			Domain:         "example.com",
			MaxSize:        10485760,
			TLSEnable:      false,
			CertFile:       "",
			KeyFile:        "",
			MaxMessageSize: "10MB",
		},
	}

	mailCore := &mockMailCore{}
	server, err := protocol.NewSMTPServer(cfg, mailCore)
	if err != nil {
		t.Fatalf("创建SMTP服务器失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Start(ctx)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("启动SMTP服务器失败: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:1025")
	if err != nil {
		t.Fatalf("连接到SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("读取欢迎消息失败: %v", err)
	}

	fmt.Fprintf(conn, "EHLO localhost\r\n")
	resp, err := readResponse(conn)
	if err != nil || !strings.HasPrefix(resp, "250-") {
		t.Errorf("EHLO命令测试失败，期望响应以'250-'开头，实际得到: %q, 错误: %v", resp, err)
	}

	fmt.Fprintf(conn, "MAIL FROM:<test@example.com>\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "250 OK\r\n" {
		t.Errorf("MAIL FROM命令测试失败，期望响应'250 OK'，实际得到: %q, 错误: %v", resp, err)
	}

	fmt.Fprintf(conn, "RCPT TO:<recipient@example.com>\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "250 OK\r\n" {
		t.Errorf("RCPT TO命令测试失败，期望响应'250 OK'，实际得到: %q, 错误: %v", resp, err)
	}

	fmt.Fprintf(conn, "DATA\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "354 End data with <CR><LF>.<CR><LF>\r\n" {
		t.Errorf("DATA命令测试失败，期望响应'354 End data with <CR><LF>.<CR><LF>'，实际得到: %q, 错误: %v", resp, err)
	}

	data := "Subject: Test\r\n\r\nThis is a test email.\r\n.\r\n"
	fmt.Fprintf(conn, "%s", data)
	resp, err = readResponse(conn)
	if err != nil || resp != "250 OK: Message accepted\r\n" {
		t.Errorf("DATA结束测试失败，期望响应'250 OK: Message accepted'，实际得到: %q, 错误: %v", resp, err)
	}

	if len(mailCore.storedEmails) != 1 {
		t.Errorf("期望存储1封邮件，实际存储了%d封", len(mailCore.storedEmails))
	} else {
		email := mailCore.storedEmails[0]
		if email.from != "test@example.com" {
			t.Errorf("期望发件人是test@example.com，实际是%s", email.from)
		}
		if len(email.recipients) != 1 || email.recipients[0] != "recipient@example.com" {
			t.Errorf("期望收件人是recipient@example.com，实际是%v", email.recipients)
		}
		if !strings.Contains(email.data, "Subject: Test\r\n\r\nThis is a test email.") {
			t.Errorf("邮件内容不匹配，实际内容: %s", email.data)
		}
	}

	fmt.Fprintf(conn, "QUIT\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "221 Bye\r\n" {
		t.Errorf("QUIT命令测试失败，期望响应'221 Bye'，实际得到: %q, 错误: %v", resp, err)
	}
}

func TestInvalidCommands(t *testing.T) {
	cfg := &config.Config{
		SMTP: config.SMTPConfig{
			Port: 1025,
		},
	}

	mailCore := &mockMailCore{}
	server, err := protocol.NewSMTPServer(cfg, mailCore)
	if err != nil {
		t.Fatalf("创建SMTP服务器失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Start(ctx)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("启动SMTP服务器失败: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:1025")
	if err != nil {
		t.Fatalf("连接到SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	testCases := []struct {
		cmd      string
		expected string
	}{
		{"INVALID", "500 Unknown command\r\n"},
		{"MAIL", "501 Syntax error in parameters or arguments\r\n"},
		{"RCPT", "501 Syntax error in parameters or arguments\r\n"},
	}

	for _, tc := range testCases {
		fmt.Fprintf(conn, "%s\r\n", tc.cmd)
		resp, err := readResponse(conn)
		if err != nil || resp != tc.expected {
			t.Errorf("命令%s测试失败，期望响应%q，实际得到: %q, 错误: %v", tc.cmd, tc.expected, resp, err)
		}
	}
}

func readResponse(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
