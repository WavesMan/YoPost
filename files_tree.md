# 完整的邮件服务器项目框架与目录结构

以下是基于Go语言开发的完整邮件服务器项目框架与目录结构：

## 项目根目录结构

```
mail-server/
├── cmd/                     # 可执行程序入口
├── configs/                 # 配置文件模板或默认配置
├── deployments/             # 部署相关文件
├── docs/                    # 文档
├── internal/                # 私有应用程序代码
├── migrations/              # 数据库迁移文件
├── pkg/                     # 可被外部导入的库代码
├── scripts/                 # 各种构建、安装、分析等脚本
├── tests/                   # 测试文件
├── third_party/             # 第三方工具、fork的代码等
├── web/                     # 前端Web应用代码
├── go.mod                   # Go模块定义
├── go.sum                   # Go模块校验
├── Makefile                 # 项目构建管理
└── README.md                # 项目说明文档
```

## 详细目录结构

### 1. cmd/ - 可执行程序入口

```
cmd/
├── server/                  # 主服务入口
│   ├── main.go              # 启动SMTP/IMAP/POP3/HTTP服务
│   └── config.go            # 服务配置加载
└── admin/                   # 管理工具
    ├── main.go              # 管理命令行工具入口
    ├── user.go              # 用户管理命令
    └── domain.go            # 域名管理命令
```

### 2. configs/ - 配置文件

```
configs/
├── server.yaml.example      # 服务配置示例
├── db.yaml.example          # 数据库配置示例
└── auth.yaml.example        # 认证配置示例
```

### 3. deployments/ - 部署相关

```
deployments/
├── docker/                  # Docker相关文件
│   ├── Dockerfile           # 后端Dockerfile
│   ├── Dockerfile.web       # 前端Dockerfile
│   └── docker-compose.yml   # 完整服务编排
└── kubernetes/              # Kubernetes部署文件
    ├── deployment.yaml
    ├── service.yaml
    └── ingress.yaml
```

### 4. internal/ - 核心应用代码

```
internal/
├── app/                     # 应用层
│   ├── server.go            # 应用主服务
│   └── shutdown.go          # 优雅关闭处理
├── auth/                    # 认证模块
│   ├── authenticator.go     # 认证接口
│   ├── database/            # 数据库认证
│   └── ldap/                # LDAP认证
├── mail/                    # 邮件处理核心
│   ├── delivery/            # 邮件投递
│   ├── storage/             # 邮件存储
│   └── processing/          # 邮件处理(过滤、分类等)
├── protocol/                # 协议实现
│   ├── smtp/                # SMTP服务
│   ├── imap/                # IMAP服务
│   └── pop3/                # POP3服务
├── api/                     # HTTP API
│   ├── v1/                  # API v1版本
│   │   ├── mail.go          # 邮件相关API
│   │   ├── account.go       # 账户相关API
│   │   └── admin.go         # 管理API
│   ├── middleware/          # API中间件
│   └── router.go            # 路由配置
├── repository/              # 数据访问层
│   ├── account.go           # 账户数据访问
│   ├── mail.go              # 邮件数据访问
│   └── db.go                # 数据库连接
├── models/                  # 数据模型
│   ├── account.go           # 账户模型
│   ├── mail.go              # 邮件模型
│   └── domain.go            # 域名模型
├── service/                 # 业务服务
│   ├── account.go           # 账户服务
│   ├── mail.go              # 邮件服务
│   └── domain.go            # 域名服务
├── webapp/                  # Web应用服务
│   ├── static.go            # 静态文件服务
│   └── websocket.go         # WebSocket服务
└── utils/                   # 工具函数
    ├── crypto.go            # 加密相关
    ├── validator.go         # 验证工具
    └── logger.go            # 日志工具
```

### 5. pkg/ - 可复用库代码

```
pkg/
├── imaputil/                # IMAP工具库
├── smtputil/                # SMTP工具库
├── mailparse/               # 邮件解析库
├── webmail/                 # Web邮件客户端库
└── auth/                    # 认证库
```

### 6. web/ - 前端Web应用

```
web/
├── public/                  # 静态文件
│   ├── index.html           # 主HTML文件
│   ├── favicon.ico          # 网站图标
│   └── assets/              # 静态资源
│       ├── images/          # 图片资源
│       └── styles/          # 全局样式
└── src/                     # 源代码
    ├── components/          # 可复用组件
    │   ├── layout/          # 布局组件
    │   ├── mail/            # 邮件相关组件
    │   └── ui/              # UI组件
    ├── pages/               # 页面组件
    │   ├── Auth/            # 认证相关页面
    │   ├── Mail/            # 邮件相关页面
    │   ├── Contacts/        # 联系人页面
    │   ├── Settings/        # 设置页面
    │   └── Admin/           # 管理页面
    ├── services/            # API服务
    │   ├── api.js           # API基础配置
    │   ├── auth.js          # 认证服务
    │   ├── mail.js          # 邮件服务
    │   └── contacts.js      # 联系人服务
    ├── store/               # 状态管理
    │   ├── index.js         # store主文件
    │   ├── auth.js          # 认证状态
    │   └── mail.js          # 邮件状态
    ├── styles/              # 样式文件
    │   ├── theme.js         # 主题配置
    │   └── global.css       # 全局样式
    ├── utils/               # 工具函数
    ├── App.js               # 主应用组件
    ├── index.js             # 入口文件
    └── routes.js            # 路由配置
```

### 7. tests/ - 测试目录

```
tests/
├── unit/                    # 单元测试
│   ├── auth/                # 认证测试
│   ├── mail/                # 邮件测试
│   └── protocol/            # 协议测试
├── integration/             # 集成测试
│   ├── api/                 # API测试
│   └── mailflow/            # 邮件流测试
└── e2e/                     # 端到端测试
    ├── web/                 # 前端测试
    └── cli/                 # CLI测试
```

### 8. 其他重要文件

```
├── .gitignore              # Git忽略规则
├── .dockerignore           # Docker忽略规则
├── Makefile                # 项目构建命令
├── go.mod                  # Go模块定义
├── go.sum                  # Go模块校验
├── package.json            # 前端依赖管理
├── package-lock.json       # 前端依赖锁文件
└── README.md               # 项目说明文档
```

## 关键文件说明

1. **cmd/server/main.go** - 主服务入口点
```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yourdomain.com/mail-server/internal/app"
	"yourdomain.com/mail-server/internal/config"
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
```

2. **internal/app/server.go** - 应用主服务
```go
package app

import (
	"context"
	"sync"

	"yourdomain.com/mail-server/internal/api"
	"yourdomain.com/mail-server/internal/mail"
	"yourdomain.com/mail-server/internal/protocol"
	"yourdomain.com/mail-server/internal/webapp"
)

type Server struct {
	cfg        *config.Config
	httpServer *api.Server
	smtpServer *protocol.SMTPServer
	imapServer *protocol.IMAPServer
	pop3Server *protocol.POP3Server
	mailCore   *mail.Core
}

func New(cfg *config.Config) (*Server, error) {
	// 初始化邮件核心
	mailCore, err := mail.NewCore(cfg)
	if err != nil {
		return nil, err
	}

	// 初始化API服务器
	httpServer, err := api.NewServer(cfg, mailCore)
	if err != nil {
		return nil, err
	}

	// 初始化协议服务器
	smtpServer := protocol.NewSMTPServer(cfg, mailCore)
	imapServer := protocol.NewIMAPServer(cfg, mailCore)
	pop3Server := protocol.NewPOP3Server(cfg, mailCore)

	return &Server{
		cfg:        cfg,
		httpServer: httpServer,
		smtpServer: smtpServer,
		imapServer: imapServer,
		pop3Server: pop3Server,
		mailCore:   mailCore,
	}, nil
}

func (s *Server) Start() error {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := s.httpServer.Start(); err != nil {
			// 处理错误
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.smtpServer.Start(); err != nil {
			// 处理错误
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.imapServer.Start(); err != nil {
			// 处理错误
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.pop3Server.Start(); err != nil {
			// 处理错误
		}
	}()

	wg.Wait()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	// 实现优雅关闭逻辑
	return nil
}
```

这个目录结构提供了清晰的代码组织方式，便于团队协作和长期维护。您可以根据实际项目需求进行调整，例如添加监控、日志等模块。需要我详细说明任何特定部分的实现细节吗？