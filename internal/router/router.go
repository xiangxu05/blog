package router

import (
	"blog/internal/endpoint"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 设置Web路由和静态文件服务
	endpoint.SetupWebRoutes(r)

	api := r.Group("/api")

	// 用户模块
	user := api.Group("/users")
	user.POST("/register", endpoint.RegisterHandler)
	user.POST("/login", endpoint.LoginHandler)
	user.Use(middleware.UserAuth())
	{
		user.GET("/logout", endpoint.LogoutHandler)
		user.GET("/get", endpoint.GetUserHandler)
		user.PUT("/update", endpoint.UpdateUserHandler)
		user.DELETE("/delete", endpoint.DeleteUserHandler)
	}

	// 文章模块
	articles := api.Group("/articles")
	// 公开访问的文章接口
	articles.GET("", endpoint.GetArticleListHandler)                              // 获取文章列表
	articles.GET("/:id", endpoint.GetArticleHandler)                              // 获取文章详情
	articles.GET("/search", endpoint.SearchArticlesHandler)                       // 搜索文章
	articles.GET("/popular", endpoint.GetPopularArticlesHandler)                  // 获取热门文章
	articles.GET("/user/:user_id", endpoint.GetUserArticlesHandler)               // 获取用户文章
	articles.GET("/category/:category_id", endpoint.GetArticlesByCategoryHandler) // 获取分类文章
	articles.GET("/categories", endpoint.GetCategoriesHandler)                    // 获取所有分类

	// 需要登录和权限的文章接口
	articles.Use(middleware.UserAuth(), middleware.RoleAuth())
	{
		articles.POST("", endpoint.CreateArticleHandler)       // 创建文章
		articles.PUT("/:id", endpoint.UpdateArticleHandler)    // 更新文章
		articles.DELETE("/:id", endpoint.DeleteArticleHandler) // 删除文章
		articles.POST("/upload-image", endpoint.UploadImageHandler) // 上传图片
	}

	return r
}
