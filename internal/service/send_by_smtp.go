package service

import (
	"YoPost/internal/config"
	"YoPost/internal/mail/core"
	"log"
)

// SendTestEmail 发送测试邮件
// to: 收件人邮箱地址
// subject: 邮件主题
// body: 邮件内容
func SendTestEmail(to string, subject string, body string) error {
	// 加载测试配置
	testCfg, err := config.LoadTestConfig()
	if err != nil {
		log.Printf("ERROR: Failed to load test config - %v", err)
		return err
	}

	// 初始化邮件服务器配置
	if err := core.InitMailServer(); err != nil {
		log.Printf("ERROR: Failed to init mail server - %v", err)
		return err
	}

	// 获取第一个测试用户
	user := testCfg.Userinfo[0]

	// 构建邮件内容
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// 调用核心邮件发送功能
	err = core.TLSstatus(user.Email, []string{to}, msg, user.Username, user.Password)
	if err != nil {
		log.Printf("ERROR: Failed to send test email - %v", err)
		return err
	}

	log.Printf("INFO: Test email sent successfully to %s", to)
	return nil
}

// 在文件末尾添加:
func main() {
	err := SendTestEmail("testuser@yopost.com", "", "")
	if err != nil {
		log.Fatalf("Failed to send test email: %v", err)
	}
	log.Println("Test email sent successfully")
}
