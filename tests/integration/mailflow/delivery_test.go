// TestMailDeliveryFlow 测试邮件投递流程的端到端集成测试
// 1. 启动PostgreSQL和Mailhog测试容器
// 2. 加载测试配置并初始化应用
// 3. 发送测试邮件并验证Mailhog中收到的邮件
// 该测试验证了从邮件发送到接收的完整流程
package mailflow_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/YoPost/internal/app"
	"github.com/YoPost/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/smtp"
	"encoding/json"
)

func startPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			Env:          map[string]string{"POSTGRES_PASSWORD": "password"},
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
	}
	return testcontainers.GenericContainer(ctx, req)
}

func startMailhogContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mailhog/mailhog",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		},
	}
	return testcontainers.GenericContainer(ctx, req)
}

func loadTestConfig(pg, mailhog testcontainers.Container) *config.Config {
	pgHost, _ := pg.Host(context.Background())
	mailhogPort, _ := mailhog.MappedPort(context.Background(), "1025")

	return &config.Config{
		DB: config.DBConfig{
			Host: pgHost,
			Port: 5432,
		},
		SMTP: config.SMTPConfig{
			Port: mailhogPort.Int(),
		},
	}
}

func sendTestEmail(t *testing.T, from, to string) {
    // 实现实际的SMTP邮件发送逻辑
    client, err := smtp.Dial(fmt.Sprintf("localhost:%d", cfg.SMTP.Port))
    if err != nil {
        t.Fatalf("SMTP连接失败: %v", err)
    }
    defer client.Close()
    
    if err := client.Mail(from); err != nil {
        t.Fatalf("MAIL FROM命令失败: %v", err)
    }
    if err := client.Rcpt(to); err != nil {
        t.Fatalf("RCPT TO命令失败: %v", err)
    }
    
    wc, err := client.Data()
    if err != nil {
        t.Fatalf("DATA命令失败: %v", err)
    }
    defer wc.Close()
    
    msg := []byte("To: " + to + "\r\nSubject: 测试邮件\r\n\r\n这是集成测试邮件内容")
    if _, err := wc.Write(msg); err != nil {
        t.Fatalf("邮件内容写入失败: %v", err)
    }
}

func getMailhogMessage(t *testing.T) struct{ From string } {
    // 实现Mailhog API查询逻辑
    var result struct {
        Items []struct {
            From struct {
                Mailbox string `json:"Mailbox"`
                Domain  string `json:"Domain"`
            } `json:"From"`
        } `json:"items"`
    }
    
    mailhogPort, _ := mailhogContainer.MappedPort(context.Background(), "8025")
    resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v2/messages", mailhogPort.Port()))
    if err != nil {
        t.Fatalf("查询Mailhog API失败: %v", err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    if err := json.Unmarshal(body, &result); err != nil {
        t.Fatalf("解析Mailhog响应失败: %v", err)
    }
    
    if len(result.Items) > 0 {
        return struct{ From string }{
            From: result.Items[0].From.Mailbox + "@" + result.Items[0].From.Domain,
        }
    }
    return struct{ From string }{}
}

func TestMailDeliveryFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pgContainer, err := startPostgresContainer(ctx)
	assert.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	mailhogContainer, err := startMailhogContainer(ctx)
	assert.NoError(t, err)
	defer mailhogContainer.Terminate(ctx)

	cfg := loadTestConfig(pgContainer, mailhogContainer)
	app, err := app.New(cfg)
	assert.NoError(t, err)

	go func() {
		assert.NoError(t, app.Start(ctx))
	}()
	defer app.Shutdown(ctx)

	sendTestEmail(t, "sender@test.com", "recipient@test.com")
	msg := getMailhogMessage(t)
	assert.Equal(t, "sender@test.com", msg.From)
}
