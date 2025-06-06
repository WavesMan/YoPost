package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动邮件服务",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("启动邮件服务...")
		// 启动SMTP服务
		cfg := &config.Config{
			SMTP: config.SMTPConfig{
				Port: 2525,
			},
			IMAP: config.IMAPConfig{
				Port: 143,
			},
			POP3: config.POP3Config{
				Port: 110,
			},
		}

		mailCore, err := mail.NewCore(cfg)
		if err != nil {
			fmt.Printf("创建邮件核心失败: %v\n", err)
			os.Exit(1)
		}

		// 检查端口是否可用
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

		// 启动所有邮件服务
		servers := []struct {
			name    string
			startFn func() error
		}{
			{
				name: "SMTP",
				startFn: func() error {
					fmt.Printf("启动SMTP服务(端口%d)...\n", cfg.SMTP.Port)
					return protocol.NewSMTPServer(cfg, mailCore).Start()
				},
			},
			{
				name: "IMAP",
				startFn: func() error {
					fmt.Printf("启动IMAP服务(端口%d)...\n", cfg.IMAP.Port)
					return protocol.NewIMAPServer(cfg, mailCore).Start()
				},
			},
			{
				name: "POP3",
				startFn: func() error {
					fmt.Printf("启动POP3服务(端口%d)...\n", cfg.POP3.Port)
					return protocol.NewPOP3Server(cfg, mailCore).Start()
				},
			},
		}

		errCh := make(chan error, len(servers))
		for _, srv := range servers {
			go func(s struct {
				name    string
				startFn func() error
			}) {
				if err := s.startFn(); err != nil {
					errCh <- fmt.Errorf("%s服务启动失败: %v", s.name, err)
				}
			}(srv)
		}

		// 等待第一个错误或全部成功
		select {
		case err := <-errCh:
			fmt.Println(err)
			os.Exit(1)
		default:
			fmt.Println("所有邮件服务已成功启动")
		}

		// 阻塞主线程
		select {}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止邮件服务",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("邮件服务停止中...")
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
