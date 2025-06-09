package api

import (
	"log"
	"net/http"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
    "fmt"
    "strings"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "context"
)

type Server struct {
	cfg      *config.Config
	mailCore mail.Core
}

func NewServer(cfg *config.Config, mailCore mail.Core) (*Server, error) {
	return &Server{
		cfg:      cfg,
		mailCore: mailCore,
	}, nil
}

func (s *Server) Start() error {
	// 初始化路由器
	mux := http.NewServeMux()

	// 挂载邮件相关API路由
	mailHandler := mail.NewMailHandler(s.mailCore)
	mux.Handle("/api/mail/", securityHeadersMiddleware(
		authMiddleware(
			http.StripPrefix("/api/mail", mailHandler),
			s.cfg.Auth.SecretKey,
		),
	))

	// 添加服务启动日志
	log.Printf("API服务启动中，监听地址: %s", s.cfg.Server.ListenAddr)
	log.Println("已注册路由:")
	log.Println("- POST /api/mail/smtp/send")

	// 启动HTTPS服务
	if s.cfg.SMTP.TLSEnable { // 修改字段引用路径
		log.Printf("Starting HTTPS server with TLS cert:%s key:%s", 
			s.cfg.SMTP.CertFile, // 修改字段引用路径
			s.cfg.SMTP.KeyFile)  // 修改字段引用路径
		return http.ListenAndServeTLS(
			s.cfg.Server.ListenAddr,
			s.cfg.SMTP.CertFile, // 修改字段引用路径
			s.cfg.SMTP.KeyFile,  // 修改字段引用路径
			mux,
		)
	}

	// 启动HTTP服务器
	return http.ListenAndServe(s.cfg.Server.ListenAddr, mux)
}

// 新增JWT认证中间件
func authMiddleware(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 跳过预检请求
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		// 解析Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			http.Error(w, `{"error": "Invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		// 验证JWT签名
		token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		// 处理验证错误
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// 验证令牌Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// 检查令牌有效期
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					http.Error(w, `{"error": "Token expired"}`, http.StatusUnauthorized)
					return
				}
			}

			// 注入用户信息到上下文
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, `{"error": "Invalid token claims"}`, http.StatusUnauthorized)
		}
	})
}

// 新增安全头中间件
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
