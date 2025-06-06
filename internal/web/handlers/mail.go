package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/YoPost/internal/mail"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	mailCore  mail.Core
	templates *template.Template
}

func NewMailHandler(mailCore mail.Core) (*MailHandler, error) {
	// 加载模板
	tmpl, err := template.ParseGlob(filepath.Join("internal", "web", "templates", "*.html"))
	if err != nil {
		return nil, err
	}

	return &MailHandler{
		mailCore:  mailCore,
		templates: tmpl,
	}, nil
}

func (h *MailHandler) RegisterRoutes(router *gin.Engine) {
	// 静态文件路由
	router.Static("/static", filepath.Join("internal", "web", "static"))

	// 邮件列表路由
	router.GET("/mail", h.mailListHandler)

	// 邮件详情路由
	router.GET("/mail/:id", h.mailDetailHandler)

	// 回复邮件路由
	router.GET("/mail/reply/:id", h.replyHandler)

	// 转发邮件路由
	router.GET("/mail/forward/:id", h.forwardHandler)
}

func (h *MailHandler) mailListHandler(c *gin.Context) {
	// 获取邮件列表数据
	emails, err := h.mailCore.GetEmails() // 需要扩展mail.Core接口
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 渲染模板
	c.HTML(http.StatusOK, "base.html", gin.H{
		"NavItems": []NavItem{
			{Name: "Inbox", Icon: "inbox", Count: 5, IsActive: true},
			{Name: "Sent", Icon: "send", Count: 0},
			{Name: "Drafts", Icon: "drafts", Count: 2},
		},
		"Emails": emails,
	})
}

func (h *MailHandler) mailDetailHandler(c *gin.Context) {
	// 获取邮件详情
	email, err := h.mailCore.GetEmail(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// 渲染模板
	c.HTML(http.StatusOK, "base.html", gin.H{
		"NavItems": []NavItem{
			{Name: "Inbox", Icon: "inbox", Count: 5},
			{Name: "Sent", Icon: "send", Count: 0},
			{Name: "Drafts", Icon: "drafts", Count: 2},
		},
		"Email": email,
	})
}

// 辅助结构体
type NavItem struct {
	Name     string
	Icon     string
	Count    int
	IsActive bool
	HTMXGet  string
}

func (h *MailHandler) replyHandler(c *gin.Context) {
	email, err := h.mailCore.GetEmail(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"NavItems": []NavItem{
			{Name: "Inbox", Icon: "inbox", Count: 5},
			{Name: "Sent", Icon: "send", Count: 0},
			{Name: "Drafts", Icon: "drafts", Count: 2},
		},
		"ReplyEmail": email,
		"IsReply":    true,
	})
}

func (h *MailHandler) forwardHandler(c *gin.Context) {
	email, err := h.mailCore.GetEmail(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"NavItems": []NavItem{
			{Name: "Inbox", Icon: "inbox", Count: 5},
			{Name: "Sent", Icon: "send", Count: 0},
			{Name: "Drafts", Icon: "drafts", Count: 2},
		},
		"ForwardEmail": email,
		"IsForward":    true,
	})
}
