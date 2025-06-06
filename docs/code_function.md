# YoPost 代码功能文档

## 1. 命令行模块 (cmd)

### 1.1 主服务 (server/main.go)

```go
// 主程序入口
func main() {
    // 初始化配置
    // 创建应用实例
    // 启动服务并监听中断信号
}
```

功能说明：
- 初始化应用配置
- 创建服务实例
- 启动HTTP和邮件协议服务
- 优雅关闭处理（30秒超时）

### 1.2 邮件服务命令行 (yomail/main.go)

```go
// 主命令结构
var rootCmd = &cobra.Command{
    Use:   "yop",
    Short: "yop - 邮件服务器控制程序",
    Long:  "YoPost 是一个完整的邮件服务器解决方案\n支持SMTP/IMAP/POP3协议",
}

// 子命令实现:
// - startCmd: 启动SMTP/IMAP/POP3服务
// - stopCmd: 停止邮件服务
// - configCmd: 配置管理
// - statusCmd: 服务状态检查
// - versionCmd: 版本信息
```

功能说明：
- 基于cobra框架的命令行接口
- 端口占用检查(isPortInUse函数)
- 并发启动多个邮件服务
- 服务状态监控
- 支持以下协议端口:
  * SMTP: 2525
  * IMAP: 143 
  * POP3: 110

## 2. 内部模块 (internal)

### 2.1 API服务 (api/server.go)

```go
// Server 提供API服务核心功能
type Server struct {
    cfg      *config.Config  // 应用配置(端口、超时等)
    mailCore mail.Core      // 邮件核心服务接口
}

// 方法:
// - NewServer(cfg, mailCore) (*Server, error)
// - Start() error (待实现)
```

功能说明：
- 计划提供RESTful API接口
- 强依赖邮件核心服务
- 配置项包括：
  * API监听端口
  * 请求超时时间
  * 认证配置
- 待实现功能：
  * 用户认证API
  * 邮件管理API
  * 服务状态API

### 2.2 应用核心 (app/server.go)

```go
// Server 主服务结构
type Server struct {
    httpServer *http.Server
    smtpServer *protocol.SMTPServer
    imapServer *protocol.IMAPServer
    pop3Server *protocol.POP3Server
    mailCore   mail.Core
}

// 方法:
// - New() *Server
// - Start() error
// - Shutdown(ctx) error (当前为空实现)
```

### 2.3 配置管理 (config/config.go)

```go
// Config 应用配置
type Config struct {
    Server   ServerConfig   // 服务器通用配置
    Database DatabaseConfig // 数据库配置
    Auth     AuthConfig     // 认证配置
    SMTP     SMTPConfig     // SMTP协议配置
    IMAP     IMAPConfig     // IMAP协议配置
    POP3     POP3Config     // POP3协议配置
}

// 方法:
// - Load() (*Config, error) (当前为空实现)
```

### 2.4 邮件核心 (mail/core.go)

```go
// Core 接口
type Core interface {
    ValidateUser(email string) bool  // 当前简单返回true
    GetConfig() *Config
    StoreEmail(from string, to []string, content string) error
}

// coreImpl 实现
type coreImpl struct {
    cfg *Config
}

// StoreEmail 实现细节：
// - 输入验证：检查发件人、收件人列表和内容是否为空
// - 存储路径：使用系统临时目录/yopost_emails/
// - 文件格式：生成唯一ID作为文件名，保存为.eml格式
// - 文件内容：包含From/To头和邮件内容
// - 错误处理：返回具体错误信息包括：
//   * 输入验证错误
//   * 目录创建失败
//   * 文件创建失败
//   * 内容写入失败
```

## 3. 协议实现 (internal/protocol)

### 3.1 IMAP协议 (imap.go)

```go
// IMAPServer 结构
type IMAPServer struct {
    cfg      *config.Config  // 应用配置
    mailCore mail.Core       // 邮件核心服务
    listener net.Listener    // 网络监听器
}

// 方法:
// - NewIMAPServer(cfg, mailCore) *IMAPServer
// - Start() error
// - GetListener() net.Listener
// - handleConnection(conn net.Conn) (实现基本IMAP协议交互)
```

### 3.2 POP3协议 (pop3.go)

```go
// POP3Server 结构
type POP3Server struct {
    cfg      *config.Config  // 应用配置
    mailCore mail.Core       // 邮件核心服务
    listener net.Listener    // 网络监听器
}

// 方法:
// - NewPOP3Server(cfg, mailCore) *POP3Server
// - Start() error
// - GetListener() net.Listener
// - handleConnection(conn net.Conn) (实现基本POP3协议交互)
```

### 3.3 SMTP协议 (smtp.go)

```go
// SMTPServer 结构
type SMTPServer struct {
    cfg         *config.Config  // 应用配置
    mailCore    mail.Core      // 邮件核心服务
    currentFrom string         // 当前会话发件人
    currentTo   []string       // 当前会话收件人列表
    listener    net.Listener   // 网络监听器
}

// 方法:
// - NewSMTPServer(cfg, mailCore) *SMTPServer
// - Start() error
// - GetListener() net.Listener
// - HandleCommand(conn net.Conn, cmd string) error (处理SMTP命令)
// - handleConnection(conn net.Conn) (完整SMTP协议实现)
```

协议支持：
- EHLO/HELO
- MAIL FROM
- RCPT TO
- DATA
- QUIT
