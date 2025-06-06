// TestPOP3Server 测试POP3服务器的基本功能，包括：
// 1. 服务器启动和监听
// 2. 客户端连接建立
// 3. 欢迎消息验证
// 4. QUIT命令处理
// 使用临时端口自动分配以避免端口冲突
// 测试完成后会自动清理服务器资源
package protocol_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	. "github.com/YoPost/internal/protocol"
)

func TestPOP3Server(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		POP3: config.POP3Config{
			Port: 0, // 让系统自动分配端口
		},
	}
	mailCore, _ := mail.NewCore(cfg)
	server := NewPOP3Server(cfg, mailCore)

	// 启动测试服务器
	serverDone := make(chan error)
	go func() {
		serverDone <- server.Start()
	}()

	// 等待服务器就绪
	select {
	case err := <-serverDone:
		t.Fatalf("POP3 server failed to start: %v", err)
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
		t.Fatalf("Failed to connect to POP3 server: %v", err)
	}
	defer conn.Close()

	// 测试欢迎消息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read welcome message: %v", err)
	}
	if !bytes.Contains(buf[:n], []byte("+OK YoPost POP3 Service Ready")) {
		t.Errorf("Unexpected welcome message: %s", string(buf[:n]))
	}

	// 测试QUIT命令
	_, err = conn.Write([]byte("QUIT\r\n"))
	if err != nil {
		t.Fatalf("Failed to write QUIT command: %v", err)
	}
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read QUIT response: %v", err)
	}
	if !bytes.Contains(buf[:n], []byte("+OK Logging out")) {
		t.Errorf("Unexpected QUIT response: %s", string(buf[:n]))
	}
}
