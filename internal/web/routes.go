package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//go:embed internal/web/dist/*
var staticFiles embed.FS

func serveFrontend(router *gin.Engine) {
	// 检查是否存在dist目录
	_, err := os.Stat("internal/web/dist")
	if os.IsNotExist(err) {
		log.Println("警告: 前端构建目录不存在，将使用嵌入式文件")
	}

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