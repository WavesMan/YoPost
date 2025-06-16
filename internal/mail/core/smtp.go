package core

import (
	"YoPost/internal/config"
	service "YoPost/services"
	"log"
	"net/smtp"
)

// MailServerConfig holds SMTP server configuration
type MailServerConfig struct {
	Host      string
	TLSPort   string
	NoTLSPort string
}

var mailServerConfig *MailServerConfig

// InitMailServer loads configuration from mailserver.yml
func InitMailServer() error {
	cfg, err := config.LoadMailServerConfig()
	if err != nil {
		return err
	}

	mailServerConfig = &MailServerConfig{
		Host:      cfg.Mailserver.Host,
		TLSPort:   cfg.Mailserver.Smtp.TlsPort,
		NoTLSPort: cfg.Mailserver.Smtp.NotlsPort,
	}

	return nil
}

// GetMailServerConfig returns the initialized mail server configuration
func GetMailServerConfig() *MailServerConfig {
	return mailServerConfig
}

func TLSstatus(from string, to []string, msg []byte, username, password string) error {
	// 检查TLS是否可用
	log.Printf("INFO: Performing fresh TLS check for %s:%s", mailServerConfig.Host, mailServerConfig.TLSPort)

	// 直接调用TLS检查服务，不再依赖缓存
	tlsEnabled := service.CheckTLSGlobal(mailServerConfig.Host, mailServerConfig.TLSPort)

	if tlsEnabled {
		// 使用TLS方式发送
		log.Printf("INFO: TLS available, attempting secure mail sending to %s:%s", mailServerConfig.Host, mailServerConfig.TLSPort)
		return SendMailWithTLS(from, to, msg, username, password)
	}
	// 使用非TLS方式发送
	log.Printf("INFO: TLS not available, attempting unsecured mail sending to %s:%s", mailServerConfig.Host, mailServerConfig.NoTLSPort)
	return SendMailWithNoTLS(from, to, msg, username, password)
}

// SendMailwithTLS 发送使用TLS加密的邮件
// from: 发件人邮箱地址
// to: 收件人邮箱地址列表
// msg: 邮件内容
// host: SMTP服务器地址
// port: SMTP服务器端口
// username: 认证用户名
// password: 认证密码
func SendMailWithTLS(from string, to []string, msg []byte, username, password string) error {
	// 使用全局配置
	host := mailServerConfig.Host
	port := mailServerConfig.TLSPort

	// 认证信息
	log.Printf("INFO: Preparing TLS mail sending to %s:%s", host, port)
	auth := smtp.PlainAuth("", username, password, host)
	log.Printf("DEBUG: Auth details - Username: %s, Host: %s", username, host)

	// 组合服务器地址
	addr := host + ":" + port
	log.Printf("INFO: Attempting secure mail sending to %s", addr)
	log.Printf("DEBUG: Message length: %d bytes", len(msg))

	// 发送邮件(使用TLS)
	c, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("ERROR: Connection failed to %s - %v", addr, err)
		return err
	}
	defer c.Close()

	if err := service.Authenticate(c, host, username, password); err != nil {
		return err
	}

	err = smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		log.Printf("ERROR: Secure mail sending failed to %s - %v", addr, err)
		return err
	}

	log.Printf("INFO: Secure mail successfully sent to %s", addr)
	return nil
}

// SendMailWithNoTLS 发送不使用TLS的邮件
// from: 发件人邮箱地址
// to: 收件人邮箱地址列表
// msg: 邮件内容
// host: SMTP服务器地址
// port: SMTP服务器端口
// username: 认证用户名
// password: 认证密码
func SendMailWithNoTLS(from string, to []string, msg []byte, username, password string) error {
	// 使用全局配置
	host := mailServerConfig.Host
	port := mailServerConfig.NoTLSPort

	// 连接SMTP服务器
	log.Printf("INFO: Preparing unsecured mail sending to %s:%s", host, port)
	addr := host + ":" + port
	log.Printf("DEBUG: Attempting connection to %s", addr)

	c, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("ERROR: Connection failed to %s - %v", addr, err)
		return err
	}
	defer c.Close()
	log.Printf("INFO: Connected to %s", addr)

	if err := service.Authenticate(c, host, username, password); err != nil {
		return err
	}

	// 设置发件人
	log.Printf("INFO: Setting sender %s", from)
	if err := c.Mail(from); err != nil {
		log.Printf("ERROR: Failed to set sender %s - %v", from, err)
		return err
	}

	// 设置收件人
	log.Printf("DEBUG: Setting recipients: %v", to)
	for _, addr := range to {
		log.Printf("INFO: Adding recipient %s", addr)
		if err := c.Rcpt(addr); err != nil {
			log.Printf("ERROR: Failed to add recipient %s - %v", addr, err)
			return err
		}
	}

	// 发送邮件内容
	log.Printf("INFO: Preparing message data")
	w, err := c.Data()
	if err != nil {
		log.Printf("ERROR: Failed to prepare message data - %v", err)
		return err
	}

	log.Printf("DEBUG: Writing message (%d bytes)", len(msg))
	_, err = w.Write(msg)
	if err != nil {
		log.Printf("ERROR: Failed to write message - %v", err)
		return err
	}

	log.Printf("DEBUG: Closing message writer")
	err = w.Close()
	if err != nil {
		log.Printf("ERROR: Failed to close message writer - %v", err)
		return err
	}

	log.Printf("INFO: Message successfully sent, quitting connection")
	return c.Quit()
}
