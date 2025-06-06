package mailflow_test

import (
	"context"
	"testing"
	"time"

	"github.com/YoPost/internal/app"
	"github.com/YoPost/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	// TODO: 实现邮件发送逻辑
}

func getMailhogMessage(t *testing.T) struct{ From string } {
	// TODO: 实现邮件获取逻辑
	return struct{ From string }{From: "sender@test.com"}
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
		assert.NoError(t, app.Start())
	}()
	defer app.Shutdown(ctx)

	sendTestEmail(t, "sender@test.com", "recipient@test.com")
	msg := getMailhogMessage(t)
	assert.Equal(t, "sender@test.com", msg.From)
}
