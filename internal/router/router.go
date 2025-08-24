package router

import (
	"blog/internal/endpoint"
	"blog/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func InitRouter() *gin.Engine {
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

	// // 文章模块
	articles := api.Group("/articles")
	// // 公开访问的文章接口
	articles.GET("/:id/:version", endpoint.GetArticleHandler) // 获取文章详情
	articles.GET("", endpoint.GetArticleListHandler)          // 获取文章列表
	// articles.GET("/search", endpoint.SearchArticlesHandler) // 搜索文章
	// articles.GET("/popular", endpoint.GetPopularArticlesHandler)                  // 获取热门文章
	// articles.GET("/user/:user_id", endpoint.GetUserArticlesHandler)               // 获取用户文章
	// articles.GET("/category/:category_id", endpoint.GetArticlesByCategoryHandler) // 获取分类文章
	// articles.GET("/categories", endpoint.GetCategoriesHandler)                    // 获取所有分类

	// // 需要登录和权限的文章接口
	articles.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		articles.POST("", endpoint.CreateArticleHandler)       // 创建文章
		articles.PUT("/:id", endpoint.UpdateArticleHandler)    // 更新文章
		articles.DELETE("/:id", endpoint.DeleteArticleHandler) // 删除文章
	}

	// 管理员专用文章接口
	admin := api.Group("/admin")
	admin.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		adminArticles := admin.Group("/articles")
		{
			adminArticles.GET("", endpoint.GetArticleListHandler)              // 获取文章列表（管理员视图）
			adminArticles.GET("/recent", endpoint.GetArticleListHandler)       // 获取最近文章
			adminArticles.GET("/:id", endpoint.GetArticleHandler)              // 获取文章详情
			adminArticles.POST("", endpoint.CreateArticleHandler)              // 创建文章
			adminArticles.PUT("/:id", endpoint.UpdateArticleHandler)           // 更新文章
			adminArticles.DELETE("/:id", endpoint.DeleteArticleHandler)        // 删除文章
			adminArticles.PUT("/:id/publish", endpoint.UpdateArticleHandler)   // 发布文章
			adminArticles.PUT("/:id/archive", endpoint.UpdateArticleHandler)   // 归档文章
			adminArticles.PUT("/batch/publish", endpoint.UpdateArticleHandler) // 批量发布
			adminArticles.PUT("/batch/archive", endpoint.UpdateArticleHandler) // 批量归档
			adminArticles.DELETE("/batch", endpoint.DeleteArticleHandler)      // 批量删除
		}
	}

	return r
}
