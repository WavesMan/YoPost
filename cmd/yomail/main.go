// YoPost 邮件服务器主程序
//
// 提供命令行接口管理邮件服务，支持以下子命令:
//   - start: 启动SMTP/IMAP/POP3服务
//   - stop: 停止邮件服务
//   - config: 配置管理
//   - status: 服务状态检查
//   - version: 显示版本信息
//
// 使用cobra框架实现命令行交互，通过检查端口占用情况确保服务正常启动。
// 主服务启动后会阻塞主线程保持运行。
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sync"

	"path/filepath"
	"syscall"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
	"github.com/spf13/cobra"
)

// 服务控制结构
type serverControl struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// 全局服务控制器
var (
	imapServer serverControl
	smtpServer serverControl
	pop3Server serverControl
)

var rootCmd = &cobra.Command{
	Use:   "yop",
	Short: "yop - 邮件服务器控制程序",
	Long: `YoPost 是一个完整的邮件服务器解决方案
支持SMTP/IMAP/POP3协议`,
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动邮件服务",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("启动邮件服务...")

		cfg := loadConfig()
		mailCore := initMailCore(cfg)
		checkPorts(cfg)

		// 初始化服务上下文
		imapServer.ctx, imapServer.cancel = context.WithCancel(context.Background())
		smtpServer.ctx, smtpServer.cancel = context.WithCancel(context.Background())
		pop3Server.ctx, pop3Server.cancel = context.WithCancel(context.Background())

		// 启动服务
		go startIMAP(cfg, mailCore)
		go startSMTP(cfg, mailCore)
		go startPOP3(cfg, mailCore)

		waitForShutdown()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止邮件服务",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("停止邮件服务...")
		if imapServer.cancel != nil {
			imapServer.cancel()
		}
		if smtpServer.cancel != nil {
			smtpServer.cancel()
		}
		if pop3Server.cancel != nil {
			pop3Server.cancel()
		}
		fmt.Println("所有服务已停止")
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "开发模式",
	Long: `启动开发环境，同时运行前端和后端服务。
前端服务会启动 yarn dev，
后端服务会启动 go run main.go`,
	Run: func(cmd *cobra.Command, args []string) {
		// 使用 WaitGroup 等待所有 goroutine 完成
		var wg sync.WaitGroup
		wg.Add(2) // 等待前端和后端两个服务

		// 启动前端开发服务器
		go func() {
			defer wg.Done()

			webCmd := exec.Command("yarn", "dev")
			webCmd.Dir = filepath.Join("web")
			webCmd.Stdout = os.Stdout
			webCmd.Stderr = os.Stderr

			log.Println("正在启动前端开发服务器...")
			if err := webCmd.Run(); err != nil {
				log.Printf("前端开发服务器启动失败: %v\n", err)
				return
			}
		}()

		// 启动后端开发服务器
		go func() {
			defer wg.Done()

			goCmd := exec.Command("go", "run", "main.go")
			goCmd.Dir = filepath.Join("cmd", "server")
			goCmd.Stdout = os.Stdout
			goCmd.Stderr = os.Stderr

			log.Println("正在启动后端开发服务器...")
			if err := goCmd.Run(); err != nil {
				log.Printf("后端开发服务器启动失败: %v\n", err)
				return
			}
		}()

		// 捕获中断信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			log.Println("接收到终止信号，正在关闭服务...")
			os.Exit(0)
		}()

		wg.Wait()
		log.Println("开发服务已全部关闭")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("配置管理功能")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "服务状态",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("服务状态检查")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("yop v0.1.0")
	},
}

func startIMAP(cfg *config.Config, mailCore mail.Core) {
	fmt.Printf("启动IMAP服务(端口%d)...\n", cfg.IMAP.Port)
	if err := protocol.NewIMAPServer(cfg, mailCore).Start(imapServer.ctx); err != nil {
		fmt.Printf("IMAP服务错误: %v\n", err)
	}
}

func startSMTP(cfg *config.Config, mailCore mail.Core) {
	fmt.Printf("启动SMTP服务(端口%d)...\n", cfg.SMTP.Port)
	if err := protocol.NewSMTPServer(cfg, mailCore).Start(smtpServer.ctx); err != nil {
		fmt.Printf("SMTP服务错误: %v\n", err)
	}
}

func startPOP3(cfg *config.Config, mailCore mail.Core) {
	fmt.Printf("启动POP3服务(端口%d)...\n", cfg.POP3.Port)
	if err := protocol.NewPOP3Server(cfg, mailCore).Start(pop3Server.ctx); err != nil {
		fmt.Printf("POP3服务错误: %v\n", err)
	}
}

func loadConfig() *config.Config {
	return &config.Config{
		SMTP: config.SMTPConfig{Port: 2525},
		IMAP: config.IMAPConfig{Port: 143},
		POP3: config.POP3Config{Port: 110},
	}
}

func initMailCore(cfg *config.Config) mail.Core {
	mailCore, err := mail.NewCore(cfg)
	if err != nil {
		fmt.Printf("创建邮件核心失败: %v\n", err)
		os.Exit(1)
	}
	return mailCore
}

func checkPorts(cfg *config.Config) {
	ports := map[string]int{
		"SMTP": cfg.SMTP.Port,
		"IMAP": cfg.IMAP.Port,
		"POP3": cfg.POP3.Port,
	}
	for name, port := range ports {
		if isPortInUse(port) {
			fmt.Printf("错误: %s端口%d已被占用\n", name, port)
			os.Exit(1)
		}
	}
}

func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\n收到停止信号，正在关闭服务...")
	stopAllServers()
}

func stopAllServers() {
	if imapServer.cancel != nil {
		imapServer.cancel()
	}
	if smtpServer.cancel != nil {
		smtpServer.cancel()
	}
	if pop3Server.cancel != nil {
		pop3Server.cancel()
	}
}

func isPortInUse(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
