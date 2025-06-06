// main 是程序的入口点，负责初始化配置、创建应用实例并启动服务。
// 它会监听系统中断信号以实现优雅关闭，超时时间为30秒。
// 如果在任何步骤中出现错误，程序将记录错误并退出。
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YoPost/internal/app"
	"github.com/YoPost/internal/config"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建应用
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// 启动服务
	go func() {
		if err := application.Start(); err != nil {
			log.Fatalf("Failed to start app: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown app: %v", err)
	}
}
