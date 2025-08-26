package endpoint

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// WebHandlerStruct 处理静态文件服务
type WebHandlerStruct struct {
	// 静态文件根目录
	StaticRoot string
}

// NewWebHandler 创建新的Web处理器
func NewWebHandler() *WebHandlerStruct {
	return &WebHandlerStruct{
		StaticRoot: "web",
	}
}

// RootHandler 根路径处理器
func RootHandler(c *gin.Context) {
	// 直接提供主页
	c.File("web/index.html")
}

// WebHandler 主要的Web处理器
func WebHandler(c *gin.Context) {
	// 提供主页HTML文件
	c.File("web/index.html")
}

// SetupWebRoutes 设置Web路由
func SetupWebRoutes(router *gin.Engine) {
	webHandler := NewWebHandler()

	// 设置静态文件服务
	webHandler.setupStaticFiles(router)

	// 设置页面路由
	webHandler.setupPageRoutes(router)

	// 设置中间件
	webHandler.SetupMiddlewares(router)
}

// setupStaticFiles 设置静态文件服务
func (w *WebHandlerStruct) setupStaticFiles(router *gin.Engine) {
	// 提供静态资源服务（CSS, JS, 图片等）
	router.Static("/assets", filepath.Join(w.StaticRoot, "assets"))
	router.Static("/templates", filepath.Join(w.StaticRoot, "templates"))

	// 提供favicon（如果存在）
	if w.FileExists("favicon.ico") {
		router.StaticFile("/favicon.ico", filepath.Join(w.StaticRoot, "favicon.ico"))
	}
}

// setupPageRoutes 设置页面路由
func (w *WebHandlerStruct) setupPageRoutes(router *gin.Engine) {
	// 主页路由
	router.GET("/", RootHandler)

	// SPA前端路由（都返回index.html，由前端路由处理）
	frontendRoutes := []string{
		"/login",
		"/register",
		"/profile",
		"/articles",
		"/articles/*path",
		"/admin",
		"/admin/*path",
	}

	for _, route := range frontendRoutes {
		router.GET(route, WebHandler)
	}

	// 处理所有非API路由的回退（SPA路由回退机制）
	router.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "API接口不存在",
				"data": nil,
			})
			return
		}

		// 其他所有请求都返回index.html，交给前端路由处理
		c.File("web/index.html")
	})
}

// HealthCheckHandler Web健康检查处理器
func (w *WebHandlerStruct) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Web服务正常运行",
		"data": gin.H{
			"status":      "healthy",
			"service":     "personal_blog_web",
			"static_root": w.StaticRoot,
		},
	})
}

// CORSMiddleware CORS中间件
func (w *WebHandlerStruct) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware 安全头中间件
func (w *WebHandlerStruct) SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止XSS攻击
		c.Header("X-XSS-Protection", "1; mode=block")

		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")

		// 内容安全策略（放宽限制以支持CDN资源）
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net https://unpkg.com https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; font-src 'self' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; img-src 'self' data: https:; connect-src 'self';")

		c.Next()
	}
}

// LoggerMiddleware 自定义日志中间件
func (w *WebHandlerStruct) LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// SetupMiddlewares 设置中间件
func (w *WebHandlerStruct) SetupMiddlewares(router *gin.Engine) {
	// 恢复中间件
	router.Use(gin.Recovery())

	// CORS中间件
	router.Use(w.CORSMiddleware())

	// 安全头中间件
	router.Use(w.SecurityHeadersMiddleware())

	// 记录中间件
}

// GetStaticPath 获取静态文件路径
func (w *WebHandlerStruct) GetStaticPath(filename string) string {
	return filepath.Join(w.StaticRoot, filename)
}

// FileExists 检查文件是否存在
func (w *WebHandlerStruct) FileExists(filename string) bool {
	path := w.GetStaticPath(filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// ServeFile 提供文件服务
func (w *WebHandlerStruct) ServeFile(c *gin.Context, filename string) {
	if !w.FileExists(filename) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文件未找到",
			"data":    nil,
		})
		return
	}

	path := w.GetStaticPath(filename)
	c.File(path)
}
