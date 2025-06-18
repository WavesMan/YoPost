# SMTP API 文档

## 发送邮件接口

`POST /api/smtp/send`

### 请求格式

API 协议说明：

1. 发送邮件 API

- 端点: POST /api/smtp/send
- 请求体:
```json
{
  "to": ["recipient@example.com"],
  "subject": "Test Subject",
  "body": "Test email content",
  "username": "your_username", // 可选 
  "password": "your_password"  // 可选
}
```

- 成功响应:
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

2. 获取 SMTP 配置 API

- 端点: GET /api/smtp/config
- 成功响应:
```json
{
  "Host": "smtp.example.com",
  "TLSPort": "465",
  "NoTLSPort": "587"
}
```