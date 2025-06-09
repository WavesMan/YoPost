// TestPOP3Server 测试POP3服务器的基本功能
// 包括服务器启动、连接建立和基本命令交互（USER/PASS/LIST/QUIT）
// 使用assert包验证服务器响应是否符合预期
// 测试用例包含：
// 1. 验证欢迎消息
// 2. 测试用户认证流程
// 3. 测试邮件列表查询
// 4. 测试正常退出流程
package protocol

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
	"github.com/stretchr/testify/assert"
)

func TestPOP3Server(t *testing.T) {
	cfg := &config.Config{
		POP3: config.POP3Config{Port: 1110},
	}
	mailCore, _ := mail.NewCore(cfg)
	server := protocol.NewPOP3Server(cfg, mailCore)

	ctx := context.Background()

	// 启动测试服务器
	go func() {
		assert.NoError(t, server.Start(ctx))
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 测试连接和基本命令
	t.Run("basic commands", func(t *testing.T) {
		conn, err := net.Dial("tcp", "localhost:1110")
		assert.NoError(t, err)
		defer conn.Close()

		// 测试欢迎消息
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "+OK YoPost POP3")

		// 测试USER/PASS命令
		_, err = conn.Write([]byte("USER test\r\n"))
		assert.NoError(t, err)
		n, err = conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "+OK")

		_, err = conn.Write([]byte("PASS 123456\r\n"))
		assert.NoError(t, err)
		n, err = conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "+OK")

		// 测试LIST命令
		_, err = conn.Write([]byte("LIST\r\n"))
		assert.NoError(t, err)
		n, err = conn.Read(buf)
		assert.NoError(t, err)
		assert.Contains(t, string(buf[:n]), "+OK 1 messages")

		// 测试QUIT命令
		_, err = conn.Write([]byte("QUIT\r\n"))
		assert.NoError(t, err)
	})
}
