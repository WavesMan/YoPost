// TestIMAPServer 测试IMAP服务器的基本功能
// 1. 启动IMAP服务器并使用动态端口
// 2. 测试基本命令流程：
//   - 验证欢迎消息
//   - 测试SELECT命令响应
//   - 测试FETCH命令响应
//   - 测试LOGOUT命令
//
// 3. 验证服务器能正常关闭
// 测试会检查网络连接、命令响应和服务器关闭行为
package protocol

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
	"github.com/stretchr/testify/assert"
)

func TestIMAPServer(t *testing.T) {
	cfg := &config.Config{
		IMAP: config.IMAPConfig{Port: 0}, // 使用动态端口
	}
	mailCore, _ := mail.NewCore(cfg)
	server := protocol.NewIMAPServer(cfg, mailCore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 直接启动服务器，让net.Listen自动选择端口
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(ctx)
	}()

	// 等待服务器真正开始监听
	var port int
	for i := 0; i < 10; i++ { // 最多尝试10次检查
		if ln := server.GetListener(); ln != nil {
			port = ln.Addr().(*net.TCPAddr).Port
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if port == 0 {
		t.Fatal("Server failed to start listening")
	}

	// 测试连接和基本命令
	t.Run("basic commands", func(t *testing.T) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 2*time.Second)
		assert.NoError(t, err)
		defer conn.Close()

		// 测试欢迎消息
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "OK YoPost IMAP")

		// 测试SELECT命令
		_, err = conn.Write([]byte("SELECT INBOX\r\n"))
		assert.NoError(t, err)
		n, err = conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "EXISTS")

		// 测试FETCH命令
		_, err = conn.Write([]byte("FETCH 1 BODY[]\r\n"))
		assert.NoError(t, err)
		n, err = conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "FETCH")

		// 测试LOGOUT命令
		_, err = conn.Write([]byte("LOGOUT\r\n"))
		assert.NoError(t, err)
	})

	// 关闭服务器
	cancel()

	// 检查服务器是否正常退出
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Error("Server did not shut down properly")
	}
}
