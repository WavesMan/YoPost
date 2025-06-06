// main 是程序的入口点，负责初始化配置、创建应用实例并启动服务。
// 它会监听系统中断信号以实现优雅关闭，超时时间为30秒。
// 如果在任何步骤中出现错误，程序将记录错误并退出。
package main

import (
	"context"
	"html/template"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YoPost/internal/app"
	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/web/handlers"
	"github.com/gin-gonic/gin"
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

	// 初始化Web服务
	webRouter := gin.Default()

	// 注册自定义模板函数
	webRouter.SetFuncMap(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// 加载模板
	webRouter.LoadHTMLGlob("internal/web/templates/**/*")

	// 设置静态文件路由
	webRouter.Static("/static", "internal/web/static")

	mailHandler, err := handlers.NewMailHandler(application.MailCore())
	if err != nil {
		log.Fatalf("Failed to create mail handler: %v", err)
	}
	mailHandler.RegisterRoutes(webRouter)

	// 启动邮件协议服务
	go func() {
		ctx := context.Background()
		if err := application.Start(ctx); err != nil {
			log.Fatalf("Failed to start mail services: %v", err)
		}
	}()

	// 启动Web服务
	go func() {
		addr := ":3000"
		log.Printf("Web服务监听在 %s", addr)
		if err := webRouter.Run(addr); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
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
