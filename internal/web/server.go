// Package web provides web interface functionality
//
// Contains:
// - Static file embedding (dist/*)
// - SPA routing handler
// - Web server configuration
package web

import (
	"embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var staticFiles embed.FS

func ServeFrontend(router *gin.Engine) {
	// 创建静态文件处理器
	fileServer := http.FileServer(http.FS(staticFiles))

	// 注册所有路径到SPA处理器
	router.NoRoute(func(c *gin.Context) {
		// 记录请求
		log.Printf("Serving static file: %s", c.Request.URL.Path)

		// 使用文件服务器提供资源
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
