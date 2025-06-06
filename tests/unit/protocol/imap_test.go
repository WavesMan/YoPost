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

func TestIMAPServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建测试配置
	cfg := &config.Config{
		IMAP: config.IMAPConfig{
			Port: 0, // 让系统自动分配端口
		},
	}
	mailCore, _ := mail.NewCore(cfg)
	server := NewIMAPServer(cfg, mailCore)

	// 启动测试服务器
	serverDone := make(chan error)
	go func() {
		serverDone <- server.Start(ctx)
	}()

	// 等待服务器就绪
	select {
	case err := <-serverDone:
		t.Fatalf("IMAP server failed to start: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	// 确保测试完成后关闭服务器
	t.Cleanup(func() {
		if ln := server.GetListener(); ln != nil {
			ln.Close()
		}
		// 等待服务器关闭
		time.Sleep(100 * time.Millisecond)
	})

	// 获取服务器地址
	serverAddr := server.GetListener().Addr().String()

	// 测试连接
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect to IMAP server: %v", err)
	}
	defer conn.Close()

	// 测试欢迎消息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read welcome message: %v", err)
	}
	if !bytes.Contains(buf[:n], []byte("* OK YoPost IMAP Service Ready")) {
		t.Errorf("Unexpected welcome message: %s", string(buf[:n]))
	}

	// 测试LOGOUT命令
	_, err = conn.Write([]byte("a001 LOGOUT\r\n"))
	if err != nil {
		t.Fatalf("Failed to write LOGOUT command: %v", err)
	}
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read LOGOUT response: %v", err)
	}
	if !bytes.Contains(buf[:n], []byte("* BYE Logging out")) {
		t.Errorf("Unexpected LOGOUT response: %s", string(buf[:n]))
	}
}
