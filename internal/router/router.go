package router

import (
	"blog/internal/endpoint"
	"blog/internal/middleware"
	"blog/internal/monitor"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router     *gin.Engine
	monitor    *monitor.Monitor
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许前端域名访问，不能使用*当credentials为true时
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
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

// NewServer 创建新的服务器实例
func NewServer(db *gorm.DB) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化路由
	s.setupRouter()

	// 初始化监控
	s.setupMonitor()

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
	user.POST("/login", endpoint.LoginHandler)
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
	}

	// 文章模块
	articles := api.Group("/articles")
	// 公开访问的文章接口
	articles.GET("/:id/:version", endpoint.GetArticleHandler) // 获取文章详情
	articles.GET("", endpoint.GetArticleListHandler)          // 获取文章列表
	articles.GET("/hot", endpoint.GetHotArticleListHandler)   // 获取热门文章
	// articles.GET("/search", endpoint.SearchArticlesHandler) // 搜索文章
	// articles.GET("/popular", endpoint.GetPopularArticlesHandler)                  // 获取热门文章
	// articles.GET("/user/:user_id", endpoint.GetUserArticlesHandler)               // 获取用户文章
	// articles.GET("/category/:category_id", endpoint.GetArticlesByCategoryHandler) // 获取分类文章
	// articles.GET("/categories", endpoint.GetCategoriesHandler)                    // 获取所有分类

	// 需要登录和权限的文章接口
	articles.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		articles.POST("", endpoint.CreateArticleHandler)       // 创建文章
		articles.PUT("/:id", endpoint.UpdateArticleHandler)    // 更新文章
		articles.DELETE("/:id", endpoint.DeleteArticleHandler) // 删除文章
	}

	// 分类模块
	category := api.Group("/category")
	{
		category.GET("", endpoint.GetCategoryListHandler) // 获取分类列表
	}

	// 管理员专用文章接口
	// admin := api.Group("/admin")
	// admin.Use(middleware.UserAuth(), middleware.RoleAuth())
	// {
	// 	adminArticles := admin.Group("/articles")
	// 	{
	// 		adminArticles.GET("", endpoint.GetArticleListHandler)              // 获取文章列表（管理员视图）
	// 		adminArticles.GET("/recent", endpoint.GetArticleListHandler)       // 获取最近文章
	// 		adminArticles.GET("/:id", endpoint.GetArticleHandler)              // 获取文章详情
	// 		adminArticles.POST("", endpoint.CreateArticleHandler)              // 创建文章
	// 		adminArticles.PUT("/:id", endpoint.UpdateArticleHandler)           // 更新文章
	// 		adminArticles.DELETE("/:id", endpoint.DeleteArticleHandler)        // 删除文章
	// 		adminArticles.PUT("/:id/publish", endpoint.UpdateArticleHandler)   // 发布文章
	// 		adminArticles.PUT("/:id/archive", endpoint.UpdateArticleHandler)   // 归档文章
	// 		adminArticles.PUT("/batch/publish", endpoint.UpdateArticleHandler) // 批量发布
	// 		adminArticles.PUT("/batch/archive", endpoint.UpdateArticleHandler) // 批量归档
	// 		adminArticles.DELETE("/batch", endpoint.DeleteArticleHandler)      // 批量删除
	// 	}
	// }

	s.router = r
}

// setupMonitor 设置监控
func (s *Server) setupMonitor() {
	s.monitor = monitor.NewMonitor(s.ctx)
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	// 启动监控
	s.monitor.Start()

	// 创建HTTP服务器
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	fmt.Printf("服务器启动在 %s\n", addr)

	// 启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("服务器启动失败: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	fmt.Println("正在关闭服务器...")

	// 创建关闭超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			fmt.Printf("HTTP服务器关闭失败: %v\n", err)
			return err
		}
	}

	// 取消上下文，通知所有组件停止
	s.cancel()

	fmt.Println("服务器已优雅关闭")
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
