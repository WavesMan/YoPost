package protocol

import (
	"context"
	"crypto/tls"
	"net"
	"net/smtp"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/protocol"
)

func TestSTARTTLS(t *testing.T) {
	cfg := &config.Config{
		SMTP: config.SMTPConfig{
			TLSEnable: true,
			CertFile:  "testdata/cert.pem",
			KeyFile:   "testdata/key.pem",
			Port:      587,
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
		if err := server.Start(ctx); err != nil {
			t.Logf("服务器错误: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// 连接服务器并执行STARTTLS
	conn, err := net.Dial("tcp", server.GetListener().Addr().String())
	if err != nil {
		t.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()

	// 验证STARTTLS流程
	client, err := smtp.NewClient(conn, "localhost")
	if err != nil {
		t.Fatalf("创建SMTP客户端失败: %v", err)
	}

	if ok, _ := client.Extension("STARTTLS"); !ok {
		t.Fatal("服务器不支持STARTTLS扩展")
	}

	if err := client.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
		t.Fatalf("STARTTLS命令失败: %v", err)
	}

	// 验证后续命令在加密通道中工作
	if err := client.Mail("test@example.com"); err != nil {
		t.Fatalf("加密通道MAIL FROM失败: %v", err)
	}
}
