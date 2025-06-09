package protocol

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/YoPost/internal/protocol"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// MockMailCore 完整实现 mail.Core 接口
type MockMailCore struct {
	emails map[string]*mail.Email
}

func NewMockMailCore() *MockMailCore {
	return &MockMailCore{
		emails: make(map[string]*mail.Email),
	}
}

func (m *MockMailCore) ValidateUser(email string) bool {
	return true
}

func (m *MockMailCore) GetConfig() *config.Config {
	certFile, keyFile := generateTempCertFiles()
	return &config.Config{
		SMTP: config.SMTPConfig{
			Port:      0, // 0 表示自动选择端口
			TLSEnable: true,
			CertFile:  certFile,
			KeyFile:   keyFile,
		},
	}
}

func (m *MockMailCore) StoreEmail(from string, to []string, data string) error {
	id := uuid.New().String()
	m.emails[id] = &mail.Email{
		ID:   id,
		From: from,
		To:   to,
		Body: data,
	}
	return nil
}

func (m *MockMailCore) GetEmails() ([]mail.Email, error) {
	var emails []mail.Email
	for _, email := range m.emails {
		emails = append(emails, *email)
	}
	return emails, nil
}

func (m *MockMailCore) GetEmail(id string) (*mail.Email, error) {
	if email, exists := m.emails[id]; exists {
		return email, nil
	}
	return nil, fmt.Errorf("email not found")
}

// generateTempCertFiles 生成临时证书文件
func generateTempCertFiles() (certFile, keyFile string) {
	dir := os.TempDir()
	certFile = filepath.Join(dir, "yopost_test_cert.pem")
	keyFile = filepath.Join(dir, "yopost_test_key.pem")

	// 如果证书已存在且有效，直接返回
	if _, err := tls.LoadX509KeyPair(certFile, keyFile); err == nil {
		return certFile, keyFile
	}

	// 生成新的证书
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"YoPost Test"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	certOut, _ := os.Create(certFile)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create(keyFile)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return certFile, keyFile
}

// connectWithRetry 带重试的连接函数
func connectWithRetry(addr string, maxRetry int) (net.Conn, error) {
	var conn net.Conn
	var err error

	for i := 0; i < maxRetry; i++ {
		conn, err = net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			return conn, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil, err
}

// assertSMTPResponse 断言SMTP响应
func assertSMTPResponse(t *testing.T, scanner *bufio.Scanner, expectedCode string) {
	require.True(t, scanner.Scan(), "预期收到服务器响应")
	response := scanner.Text()
	require.True(t, strings.HasPrefix(response, expectedCode),
		"预期响应码 %s, 实际收到: %s", expectedCode, response)
}

// getListenerAddr 通过反射获取监听地址
func getListenerAddr(server *protocol.SMTPServer) (string, error) {
	// 使用反射访问未导出的 listener 字段
	val := reflect.ValueOf(server).Elem()
	listenerField := val.FieldByName("listener")
	if !listenerField.IsValid() {
		return "", fmt.Errorf("listener field not found")
	}

	listener, ok := listenerField.Interface().(net.Listener)
	if !ok || listener == nil {
		return "", fmt.Errorf("invalid listener")
	}

	return listener.Addr().String(), nil
}

// TestSTARTTLS 测试STARTTLS功能
func TestSTARTTLS(t *testing.T) {
	// 1. 准备测试环境
	mailCore := NewMockMailCore()
	cfg := mailCore.GetConfig()

	// 2. 启动SMTP服务器
	server, err := protocol.NewSMTPServer(cfg, mailCore)
	require.NoError(t, err, "创建SMTP服务器失败")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverReady := make(chan struct{})
	go func() {
		close(serverReady)
		if err := server.Start(ctx); err != nil {
			t.Logf("服务器退出: %v", err)
		}
	}()

	// 等待服务器启动
	<-serverReady
	time.Sleep(100 * time.Millisecond) // 确保监听器已设置

	// 获取监听地址
	listenAddr := server.Addr()
	require.NotEmpty(t, listenAddr, "服务器监听地址为空")
	t.Logf("SMTP服务器监听地址: %s", listenAddr)

	// 3. 获取监听端口
	_, portStr, err := net.SplitHostPort(listenAddr)
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)
	t.Logf("SMTP服务器监听端口: %d", port)

	// 4. 建立连接 (带重试)
	conn, err := connectWithRetry(fmt.Sprintf("127.0.0.1:%d", port), 3)
	require.NoError(t, err, "连接服务器失败")
	defer conn.Close()

	// 5. 开始SMTP会话测试
	scanner := bufio.NewScanner(conn)
	writer := io.Writer(conn)

	// 检查欢迎消息
	assertSMTPResponse(t, scanner, "220")

	// 发送EHLO命令
	fmt.Fprintf(writer, "EHLO localhost\r\n")
	assertSMTPResponse(t, scanner, "250")

	// 发送STARTTLS命令
	fmt.Fprintf(writer, "STARTTLS\r\n")
	assertSMTPResponse(t, scanner, "220")

	// 6. 升级到TLS连接
	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true, // 测试环境跳过证书验证
		ServerName:         "localhost",
	})
	defer tlsConn.Close()

	// 需要重新创建scanner和writer
	tlsScanner := bufio.NewScanner(tlsConn)
	tlsWriter := io.Writer(tlsConn)

	// 再次发送EHLO (TLS模式下)
	fmt.Fprintf(tlsWriter, "EHLO localhost\r\n")
	assertSMTPResponse(t, tlsScanner, "250")

	// 测试MAIL FROM命令
	fmt.Fprintf(tlsWriter, "MAIL FROM:<test@yopost.com>\r\n")
	assertSMTPResponse(t, tlsScanner, "250")

	// 测试QUIT命令
	fmt.Fprintf(tlsWriter, "QUIT\r\n")
	assertSMTPResponse(t, tlsScanner, "221")
}

// TestImplicitTLS 测试隐式TLS连接
func TestImplicitTLS(t *testing.T) {
	mailCore := NewMockMailCore()
	cfg := mailCore.GetConfig()
	server, err := protocol.NewSMTPServer(cfg, mailCore)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 启动服务器
	serverReady := make(chan struct{})
	go func() {
		close(serverReady)
		if err := server.Start(ctx); err != nil {
			t.Logf("服务器退出: %v", err)
		}
	}()

	// 等待服务器启动
	<-serverReady
	time.Sleep(100 * time.Millisecond) // 确保监听器已设置

	// 获取监听地址
	listenAddr := server.Addr()
	require.NotEmpty(t, listenAddr, "服务器监听地址为空")
	t.Logf("SMTP服务器监听地址: %s", listenAddr)

	// 获取监听端口
	_, portStr, err := net.SplitHostPort(listenAddr)
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)
	// 直接建立TLS连接
	conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port), &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "localhost",
	})
	require.NoError(t, err)
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	writer := io.Writer(conn)
	// 检查欢迎消息
	assertSMTPResponse(t, scanner, "220")
	// 发送EHLO命令
	fmt.Fprintf(writer, "EHLO localhost\r\n")
	assertSMTPResponse(t, scanner, "250")
}
