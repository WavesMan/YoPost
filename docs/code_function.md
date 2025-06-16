# YoPost 代码功能文档

## 1. 后端功能

### 1.1 邮件核心

#### 1.1.1 SMTP
1. `InitMailServer` 加载 `internal/config` 配置文件
   载入 `Host` `TLSPort` `NoTLSPort`
2. `TLSstatus` 调用 `tls.go` 中的 `TLSstatus` 函数，判断是否开启TLS
3. `SendMailWithTLS` 以TLS方式发送邮件
4. `SendMailWithNoTLS` 以NoTLS方式发送邮件

### 1.2 Services 公共包
- `tls.go`：提供全局TLS状态验证功能
- `authenticate.go`：提供全局用户状态基础验证功能（暂未加密）

## 2. 前端构建 ( /web/*)
### 2.1 框架
通过现代化 Vite + React 构建，使用 JavaScript 实现逻辑处理

### 2.2 技术栈
- React：作为核心 UI 库。
- Vite：用于构建开发服务器和打包生产环境资源。
- React Router DOM (v7)：处理应用中的导航（如侧边栏不同文件夹）。
- CSS：组件级别的样式设计，未使用 CSS-in-JS 或模块化导入方式。
- ESLint：代码规范检查。

### 2.3 目录结构:
```
src/
├── components/           // 所有组件
│   ├── ComposeEmail.jsx/css    // 编写新邮件的弹窗组件
│   ├── EmailList.jsx/css       // 邮件列表展示组件
│   ├── EmailView.jsx/css       // 单封邮件预览组件
│   └── Sidebar.jsx/css         // 左侧导航菜单组件
├── App.jsx/css             // 根组件，整合所有子组件
├── main.jsx                // 入口点
└── index.css               // 全局样式
```

### 2.4 组件化实现
| 组件 | 状态 | 实现功能 |
| ---- | ---- | ---- |
| EmailList | ✅完成	| 邮件列表展示 |
| EmailView	| ✅完成	| 邮件内容展示 |
| Sidebar	| ✅完成	| 左侧导航栏 |
| ComposeEmail | 🟡开发中 | 新邮件编写 |