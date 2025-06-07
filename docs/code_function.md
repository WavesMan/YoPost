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
var startCmd = &cobra.Command{
    Use:   "start",
    Short: "启动邮件服务",
    Run: func(cmd *cobra.Command, args []string) {
        // 初始化服务上下文
        // 并发启动IMAP/SMTP/POP3服务
    }
}

var stopCmd = &cobra.Command{
    Use:   "stop", 
    Short: "停止邮件服务",
    Run: func(cmd *cobra.Command, args []string) {
        // 调用各服务的cancel函数
        // 等待所有服务停止
    }
}

// 服务控制结构
type serverControl struct {
    ctx    context.Context
    cancel context.CancelFunc
}

// 核心函数:
func startAllServers(cfg *config.Config, mailCore mail.Core) error {
    // 初始化服务上下文
    // 启动各协议服务协程
}

func stopAllServers() {
    // 调用各服务的cancel函数
}

func waitForShutdown() {
    // 监听系统信号
    // 触发优雅关闭
}
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

### 2.0 Web界面处理 (web/handlers)
```go
// MailHandler 处理邮件相关的Web界面请求
//
// 主要功能：
// - 模板加载与渲染
// - 邮件列表/详情展示
// - 邮件回复/转发处理
// - 静态资源服务
//
// 技术实现：
// - 使用html/template进行服务端渲染
// - 集成HTMX实现局部更新
// - Alpine.js处理客户端交互

type MailHandler struct {
    mailCore  mail.Core       // 邮件核心服务
    templates *template.Template // 模板集合
}

// 关键方法:
// - NewMailHandler() 初始化模板和路由
// - RegisterRoutes() 注册Web路由
// - mailListHandler() 邮件列表
// - mailDetailHandler() 邮件详情
// - replyHandler() 回复邮件
// - forwardHandler() 转发邮件
```

#### 模板系统
- 基础模板: base.html
  - 包含全局布局和静态资源
  - 支持暗黑模式切换
  - 导航菜单动态生成
- 邮件模板: mail/*.html
  - 列表视图(list.html)
  - 详情视图(detail.html)
  - 回复/转发表单(reply.html)

#### 静态资源
- CSS: /static/css/main.css
- JS: 
  - /static/htmx.min.js
  - /static/alpine.js

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

### 3.1 IMAP协议 (imap.go) - 开发进度: 65%
```go
// IMAPServer 实现了IMAP协议服务端功能
//
// 已完成功能:
// - 基础连接管理
// - 基本命令支持(LOGOUT/SELECT/FETCH/SEARCH)
// - 上下文感知的服务启停
//
// 待实现功能:
// - 邮箱状态维护(UIDVALIDITY/UIDNEXT)
// - 扩展命令支持(UIDPLUS/CONDSTORE)
// - 邮件标记操作(FLAGS/PERMANENTFLAGS)
//
// 测试覆盖率: 45%
type IMAPServer struct {
    cfg      *config.Config  // 应用配置
    mailCore mail.Core       // 邮件核心服务
    listener net.Listener    // 网络监听器
}
```

### 3.2 POP3协议 (pop3.go) - 开发进度: 70%
```go
// POP3Server 实现了POP3协议服务器
//
// 已完成功能:
// - 基础连接管理
// - 认证流程(USER/PASS)
// - 邮件列表操作(LIST/RETR/DELE)
// - 会话终止(QUIT)
//
// 待实现功能:
// - UIDL命令支持
// - TOP命令实现
// - 认证加密支持
//
// 测试覆盖率: 60%
type POP3Server struct {
    cfg      *config.Config  // 应用配置
    mailCore mail.Core       // 邮件核心服务
    listener net.Listener    // 网络监听器
}
```

### 3.3 SMTP协议 (smtp.go) - 开发进度: 85%
```go
// SMTPServer 实现了SMTP协议服务器
//
// 已完成功能:
// - 完整SMTP命令支持(EHLO/MAIL/RCPT/DATA/QUIT)
// - 邮件内容解析与存储
// - 多收件人处理
// - 错误处理与状态码返回
//
// 待实现功能:
// - STARTTLS加密支持
// - 发件人验证(SPF/DKIM)
// - 速率限制
//
// 测试覆盖率: 75%
type SMTPServer struct {
    cfg         *config.Config  // 应用配置
    mailCore    mail.Core      // 邮件核心服务
    currentFrom string         // 当前会话发件人
    currentTo   []string       // 当前会话收件人列表
    listener    net.Listener   // 网络监听器
}
```

### 3.2 POP3协议 (pop3.go)

```go
// POP3Server 实现了POP3协议服务器
//
// 主要功能包括：
// - 监听指定端口接收客户端连接
// - 处理基本的POP3命令交互
// - 提供邮件服务核心接口
//
// 使用NewPOP3Server创建实例，通过Start方法启动服务
type POP3Server struct {
    cfg      *config.Config  // 应用配置
    mailCore mail.Core       // 邮件核心服务
    listener net.Listener    // 网络监听器
}

// 方法:
// - NewPOP3Server(cfg, mailCore) *POP3Server
// - Start(ctx context.Context) error
// - GetListener() net.Listener
// - handleConnection(conn net.Conn) (实现基本POP3协议交互)

// 支持命令:
// - USER <username>
// - PASS <password>
// - LIST
// - RETR <message>
// - DELE <message>
// - QUIT
```

### 3.3 SMTP协议 (smtp.go)

```go
// SMTPServer 实现了简单的SMTP协议服务器，用于接收和处理电子邮件
//
// 主要功能包括：
// - 监听指定端口接收SMTP连接
// - 处理标准SMTP命令（EHLO/HELO、MAIL FROM、RCPT TO、DATA、QUIT等）
// - 存储接收到的邮件到邮件核心系统
//
// 结构体包含配置信息、邮件核心处理模块和当前会话状态
type SMTPServer struct {
    cfg         *config.Config  // 应用配置
    mailCore    mail.Core      // 邮件核心服务
    currentFrom string         // 当前会话发件人
    currentTo   []string       // 当前会话收件人列表
    listener    net.Listener   // 网络监听器
}

// 方法:
// - NewSMTPServer(cfg, mailCore) *SMTPServer
// - Start(ctx context.Context) error
// - GetListener() net.Listener
// - HandleCommand(conn net.Conn, cmd string) error
// - handleConnection(conn net.Conn) (完整SMTP协议实现)

// 支持命令:
// - EHLO/HELO <domain>
// - MAIL FROM:<sender>
// - RCPT TO:<recipient>
// - DATA
// - QUIT
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
