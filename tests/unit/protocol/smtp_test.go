package protocol

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/protocol"
)

func TestSMTPServer(t *testing.T) {
	mailCore := &MockMailCore{}
	// 使用动态端口配置
	cfg := &config.Config{
		SMTP: config.SMTPConfig{
			Port:      0, // 使用0让系统自动分配端口
			MaxSize:   10485760,
			TLSEnable: false,
		},
	}
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

	addr := server.GetListener().Addr().String()
	conn, err := net.Dial("tcp", addr)
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

	fmt.Fprintf(conn, "MAIL From:<test@example.com>\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "250 OK\r\n" {
		t.Errorf("MAIL From命令测试失败，期望响应'250 OK'，实际得到: %q, 错误: %v", resp, err)
	}

	fmt.Fprintf(conn, "RCPT TO:<recipient@example.com>\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "250 OK\r\n" {
		t.Errorf("RCPT TO命令测试失败，期望响应'250 OK'，实际得到: %q, 错误: %v", resp, err)
	}

	fmt.Fprintf(conn, "Data\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "354 End Data with <CR><LF>.<CR><LF>\r\n" {
		t.Errorf("Data命令测试失败，期望响应'354 End Data with <CR><LF>.<CR><LF>'，实际得到: %q, 错误: %v", resp, err)
	}

	Data := "Subject: Test\r\n\r\nThis is a test email.\r\n.\r\n"
	fmt.Fprintf(conn, "%s", Data)
	resp, err = readResponse(conn)
	if err != nil || resp != "250 OK: Message accepted\r\n" {
		t.Errorf("Data结束测试失败，期望响应'250 OK: Message accepted'，实际得到: %q, 错误: %v", resp, err)
	}

	if len(mailCore.StoredEmails) != 1 {
		t.Errorf("期望存储1封邮件，实际存储了%d封", len(mailCore.StoredEmails))
	} else {
		email := mailCore.StoredEmails[0]
		if email.From != "test@example.com" {
			t.Errorf("期望发件人是test@example.com，实际是%s", email.From)
		}
		if len(email.Recipients) != 1 || email.Recipients[0] != "recipient@example.com" {
			t.Errorf("期望收件人是recipient@example.com，实际是%v", email.Recipients)
		}
		if !strings.Contains(email.Data, "Subject: Test\r\n\r\nThis is a test email.") {
			t.Errorf("邮件内容不匹配，实际内容: %s", email.Data)
		}
	}

	fmt.Fprintf(conn, "QUIT\r\n")
	resp, err = readResponse(conn)
	if err != nil || resp != "221 Bye\r\n" {
		t.Errorf("QUIT命令测试失败，期望响应'221 Bye'，实际得到: %q, 错误: %v", resp, err)
	}
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

func TestSMTPServerPortSelection(t *testing.T) {
	testCases := []struct {
		name       string
		tlsEnabled bool
		expected   int
	}{
		{"NonTLS mode", false, 25},
		{"TLS mode", true, 465},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				Server: config.ServerConfig{Host: "127.0.0.1"},
				SMTP: config.SMTPConfig{
					TLSEnable: tc.tlsEnabled,
				},
			}

			mailCore := &MockMailCore{}
			server, err := protocol.NewSMTPServer(cfg, mailCore)
			if err != nil {
				t.Fatalf("创建SMTP服务器失败: %v", err)
			}

			// 初始化监听器
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				if err := server.Start(ctx); err != nil {
					t.Logf("服务器启动错误: %v", err)
				}
			}()
			time.Sleep(100 * time.Millisecond) // 等待服务器启动

			addr := server.GetListener().Addr().String()
			port := parsePortFromAddr(addr)
			if port != tc.expected {
				t.Errorf("期望端口%d，实际得到%d", tc.expected, port)
			}
		})
	}
}

func TestInvalidCommands(t *testing.T) {
	cfg := &config.Config{
		SMTP: config.SMTPConfig{
			Port: 1025,
		},
	}

	mailCore := &MockMailCore{}
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
