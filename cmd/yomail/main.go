package main

import (
	"fmt"
	"os"

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
		fmt.Println("启动SMTP服务...")
		// TODO: 实现IMAP/POP3服务启动
		cfg := &config.Config{
			SMTP: config.SMTPConfig{
				Port: 2525,
			},
		}

		mailCore, err := mail.NewCore(cfg)
		if err != nil {
			fmt.Printf("创建邮件核心失败: %v\n", err)
			os.Exit(1)
		}

		server := protocol.NewSMTPServer(cfg, mailCore)
		if err := server.Start(); err != nil {
			fmt.Printf("启动失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("SMTP服务已启动")
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
