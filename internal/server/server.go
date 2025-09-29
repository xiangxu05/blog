package server

import (
	"blog/internal/endpoint"
	"blog/internal/middleware"
	"blog/internal/monitor"
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiangxu05/logger/v2"
	"gorm.io/gorm"
)

type Server struct {
	router     *gin.Engine
	monitor    *monitor.Monitor
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的Origin
		origin := c.Request.Header.Get("Origin")
		// 定义允许的域名列表
		allowedOrigins := map[string]bool{
			"http://localhost:8080": true, // 本地开发环境
			"http://127.0.0.1:8080": true, // 本地开发环境
			// "https://www.yourdomain.com":  true, // 生产环境www子域名
			// "https://blog.yourdomain.com": true, // 博客子域名
		}
		// 检查请求的Origin是否在允许列表中
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		// 允许的HTTP方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// 允许凭证（Cookie、HTTP认证等）
		c.Header("Access-Control-Allow-Credentials", "true")
		// 预检请求的缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400") // 24小时
		// 如果是预检请求(OPTIONS)，直接返回204状态码
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// NewServer 创建新的服务器实例
func NewServer(db *gorm.DB) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		cancel: cancel,
		wg:     new(sync.WaitGroup),
	}

	// 初始化监控
	s.setupMonitor()

	// 初始化路由
	s.setupRouter()

	return s
}

// setupRouter 设置路由
func (s *Server) setupRouter() {
	r := gin.Default()

	// 设置Web路由和静态文件服务
	endpoint.SetupWebRoutes(r)

	// 为API路由组添加CORS中间件
	api := r.Group("/api")
	api.Use(corsMiddleware())

	api.GET("/health", endpoint.HealthHandler)
	// 用户模块
	user := api.Group("/users")
	user.POST("/register", endpoint.RegisterHandler)
	user.POST("/login", s.monitor.UpdateLastLogin(), endpoint.LoginHandler)
	// 公开访问的用户信息接口
	user.GET("/:user_id", endpoint.GetOtherUserHandler)
	user.Use(middleware.UserAuth())
	{
		user.GET("/logout", endpoint.LogoutHandler)
		user.DELETE("/delete", endpoint.DeleteUserHandler)
		user.GET("/refresh", endpoint.RefreshHandler)
		user.GET("/profile", endpoint.GetUserHandler)
		user.PUT("/profile", endpoint.UpdateUserHandler)
		user.PUT("/password", endpoint.UpdatePasswordHandler)
	}

	// 文件模块
	file := api.Group("/files")
	file.GET("/:file_id", endpoint.GetFileHandler)
	file.GET("/download/:file_id", endpoint.DownloadFileHandler)
	file.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		file.GET("/", endpoint.GetFileListHandler)
		file.POST("/upload", endpoint.UploadFileHandler)
		file.DELETE("/:file_id", endpoint.DeleteFileHandler)
		file.GET("/backup", s.monitor.UpdateLastBackup(), endpoint.BackupHandler)
	}

	// 文章模块
	articles := api.Group("/articles")
	// 公开访问的文章接口
	articles.GET("/:id/:version", s.monitor.IncViews(), endpoint.GetArticleHandler)     // 获取文章详情
	articles.GET("/:id", s.monitor.IncViews(), endpoint.GetLatestArticleHandler)        // 获取最新文章详情
	articles.GET("/:id/latest", s.monitor.IncViews(), endpoint.GetLatestArticleHandler) // 获取最新文章详情
	articles.GET("", endpoint.GetArticleListHandler)                                    // 获取文章列表
	articles.GET("/hot", endpoint.GetHotArticleListHandler)                             // 获取热门文章
	// articles.GET("/search", endpoint.SearchArticlesHandler) // 搜索文章
	// articles.GET("/popular", endpoint.GetPopularArticlesHandler)                  // 获取热门文章
	// articles.GET("/user/:user_id", endpoint.GetUserArticlesHandler)               // 获取用户文章
	// articles.GET("/category/:category_id", endpoint.GetArticlesByCategoryHandler) // 获取分类文章
	// articles.GET("/categories", endpoint.GetCategoriesHandler)                    // 获取所有分类

	// 需要登录和权限的文章接口
	articles.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		articles.POST("", s.monitor.UpdateLastUpdate(), endpoint.CreateArticleHandler) // 创建文章
		articles.PUT("/:id", endpoint.UpdateArticleHandler)                            // 更新文章
		articles.DELETE("/:id", endpoint.DeleteArticleHandler)                         // 删除文章
	}

	// 分类模块
	category := api.Group("/category")
	{
		category.GET("", endpoint.GetCategoryListHandler) // 获取分类列表
	}

	// 管理员专用文章接口
	admin := api.Group("/admin")
	admin.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		adminArticles := admin.Group("/articles")
		{
			adminArticles.GET("", endpoint.GetAdminArticleListHandler)  // 获取文章列表（管理员视图）
			adminArticles.GET("/:id", endpoint.GetArticleHandler)       // 获取文章详情
			adminArticles.POST("", endpoint.CreateArticleHandler)       // 创建文章
			adminArticles.PUT("/:id", endpoint.UpdateArticleHandler)    // 更新文章
			adminArticles.DELETE("/:id", endpoint.DeleteArticleHandler) // 删除文章
		}
	}

	// 网站信息查询
	api.GET("/website_statistics", s.monitor.GetWebsiteStatisticsHandler())
	admin.GET("/backend_statistics", middleware.RoleAuth(), s.monitor.GetBackendStatisticsHandler()) // 管理员专用
	s.router = r
}

// setupMonitor 设置监控
func (s *Server) setupMonitor() {
	s.monitor = monitor.NewMonitor(s.ctx)
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	var log = logger.GetLogger()
	// 启动监控
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.monitor.Start()
	}()

	// 创建HTTP服务器
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Infof("服务器启动在 %s", addr)

	// 启动服务器
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("服务器启动失败: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	var log = logger.GetLogger()
	// 创建关闭超时上下文
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Errorf("HTTP服务器关闭失败: %v\n", err)
			return err
		}
	}

	// 取消上下文，通知所有组件停止
	s.cancel()
	s.wg.Wait()
	log.Info("服务器已关闭")
	return nil
}

// Run 运行服务器并处理信号
func (s *Server) Run(addr string) error {
	// 启动服务器
	if err := s.Start(addr); err != nil {
		return err
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	return s.Stop()
}
