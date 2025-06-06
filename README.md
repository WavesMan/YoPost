# YoPost - 邮件服务器项目

YoPost是一个基于Go语言开发的完整邮件服务器解决方案，支持SMTP、IMAP和POP3协议。

## 功能特性

- 完整的邮件收发功能
- 多协议支持(SMTP/IMAP/POP3)
- RESTful API管理接口
- 可扩展的认证系统
- 高性能邮件处理核心

## 项目结构

```
YoPost/
├── cmd/                     # 可执行程序入口
├── configs/                 # 配置文件模板
├── deployments/             # 部署相关文件
├── docs/                    # 文档
├── internal/                # 核心应用代码
├── migrations/              # 数据库迁移
├── pkg/                     # 可复用库代码
├── scripts/                 # 构建/安装脚本
├── tests/                   # 测试文件
├── third_party/             # 第三方工具
└── web/                     # 前端Web应用
```

## 快速开始

1. 克隆项目
2. 配置环境变量
3. 启动服务:
   ```bash
   go run cmd/server/main.go
   ```

## 依赖

- Go 1.21+
- Redis (可选)
- PostgreSQL (可选)

## 许可证

MIT
