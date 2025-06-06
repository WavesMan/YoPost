// 需要添加的包导入
import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/YoPost/internal/app"
	"github.com/YoPost/internal/config"
)

// 需要实现的辅助函数
func loadTestConfig(pg, mailhog testcontainers.Container) *config.Config {
	// 获取容器IP和端口
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

func TestMailDeliveryFlow(t *testing.T) {
	// 启动测试容器(PostgreSQL + MailHog)
	ctx := context.Background()
	pgContainer := startPostgresContainer(ctx)
	mailhogContainer := startMailhogContainer(ctx)

	// 初始化测试服务
	cfg := loadTestConfig(pgContainer, mailhogContainer)
	app := app.New(cfg)
	go app.Start()

	// 发送测试邮件
	sendTestEmail(t, "sender@test.com", "recipient@test.com")

	// 验证邮件接收
	msg := getMailhogMessage(t)
	assert.Equal(t, "sender@test.com", msg.From)

	// 清理
	app.Shutdown(ctx)
}
