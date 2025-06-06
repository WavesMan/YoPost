// TestSMTPServer 测试SMTP服务器的基本功能，包括：
// 1. 服务器启动和监听
// 2. 客户端连接建立
// 3. 欢迎消息验证
// 4. EHLO命令处理
// 使用临时端口(0)进行测试，确保测试隔离性
// 包含清理逻辑确保测试后资源释放
package protocol_test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	. "github.com/YoPost/internal/protocol"
)

func TestSMTPServer(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		SMTP: config.SMTPConfig{
			Port: 0, // 让系统自动分配端口
		},
	}
	mailCore, _ := mail.NewCore(cfg)
	server := NewSMTPServer(cfg, mailCore)

	ctx := context.Background()

	// 启动测试服务器
	serverDone := make(chan error)
	go func() {
		serverDone <- server.Start(ctx)
	}()

	// 等待服务器就绪
	select {
	case err := <-serverDone:
		t.Fatalf("SMTP server failed to start: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	// 确保测试完成后关闭服务器
	t.Cleanup(func() {
		if ln := server.GetListener(); ln != nil {
			ln.Close()
		}
	})

	// 获取服务器地址
	serverAddr := server.GetListener().Addr().String()

	// 测试连接
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// 测试欢迎消息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read welcome message: %v", err)
	}
	if !bytes.Contains(buf[:n], []byte("220 YoPost SMTP Service Ready")) {
		t.Errorf("Unexpected welcome message: %s", string(buf[:n]))
	}

	// 测试EHLO命令
	_, err = conn.Write([]byte("EHLO test.example.com\r\n"))
	if err != nil {
		t.Fatalf("Failed to write EHLO command: %v", err)
	}
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read EHLO response: %v", err)
	}
	if !bytes.Contains(buf[:n], []byte("250-yop.example.com")) {
		t.Errorf("Unexpected EHLO response: %s", string(buf[:n]))
	}
}
