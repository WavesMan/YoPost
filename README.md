# YoPost - 全栈邮件服务器解决方案（早期开发中）

YoPost是一个基于Go语言开发的一体化邮件服务器，提供完整的SMTP/IMAP/POP3协议支持，以及现代化的Web管理界面。

## 核心特性

### 后端功能
- 完整的邮件服务器功能，兼容docker-mailserver
- 多协议支持：
  - SMTP (开发进度90%)
  - IMAP (开发进度65%) 
  - POP3 (开发进度70%)
- 高性能邮件处理核心
- RESTful管理API (开发中)
- 支持Docker容器化部署 (开发中)

### 前端功能
- 管理员界面：
  - [x] 用户账号管理
  - [ ] 邮件域配置
  - [ ] 服务器监控
- 用户邮箱界面：
  - [x] 邮件列表
  - [x] 邮件查看
  - [ ] 邮件回复/转发
  - [ ] 联系人管理
  - [ ] 邮件搜索

## 技术栈

### 前端
- 框架: React 19.1.0 + React Router DOM v7.6.2
- 构建工具: Vite 6.3.5
- 样式方案: CSS Modules + Flexbox布局

### 后端
- Go 1.21+
- Cobra CLI框架
- PostgreSQL/Redis (可选)
- Prometheus监控

## 项目结构

- [开发状态报告](./docs/DEV_STATUS.md) - 详细的项目进度和计划
- [代码功能文档](./docs/code_function.md) - 核心模块和实现细节
- [Wev开发文档](./docs/WEB_DEVELOPMENT.md) - Web界面开发指南
- API文档
  - [SMTP API 文档](./docs/api/SMTP-API.md)

```
YoPost/
├── cmd/             # 命令行入口
│   ├── server/      # 主服务
│   └── yomail/      # 邮件服务控制
├── internal/        # 核心实现
│   ├── api/         # REST API
│   ├── protocol/    # 邮件协议实现
│   └── webapp/      # 前端集成
├── web/
│   └── src/
│       ├── components/           // 所有组件
│       ├── App.jsx/css             // 根组件，整合所有子组件
│       ├── main.jsx                // 入口点
│       └── index.css               // 全局样式
├── tests/           # 测试代码
└── docs/            # 文档
```

## 快速开始

1. 克隆项目
   ```bash
   git clone https://github.com/yourrepo/yopost.git
   ```

2. 启动开发环境
   ```bash
   # 启动后端服务
   go run cmd/server/main.go

   # 静态资源已自动服务
   ```

3. 访问管理界面
   ```
   http://localhost:3000/admin
   ```

## 开发状态

当前版本：0.3.0 (Beta)

- 已完成：
  - 基础邮件协议实现
  - 邮件存储核心
  - 命令行控制工具

- 开发中：
  - Web管理界面
  - 用户认证系统
  - API文档生成

## 项目文档

- [开发状态报告](./docs/DEV_STATUS.md) - 详细的项目进度和计划
- [代码功能文档](./docs/code_function.md) - 核心模块和实现细节
- [Wev开发文档](./docs/WEB_DEVELOPMENT.md) - Web界面开发指南

## 贡献指南

欢迎通过Issues或Pull Requests参与贡献。请先阅读上述文档。

## 许可证

[GPL-3.0 License](License)
