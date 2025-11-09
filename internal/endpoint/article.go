package endpoint

import (
	"blog/internal/message"
	"blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiangxu05/logger/v2"
)

var log = logger.GetLogger()

// GetArticleHandler 获取文章
func GetArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleId, err := strconv.Atoi(articleStr)
	if err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	articleVersion, err := service.GetArticle(c, int32(articleId), version)
	if err != nil {
		log.Errorf("获取文章失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "文章不存在", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", articleVersion)
}

// GetLatestArticleHandler 获取最新文章
func GetLatestArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleId, err := strconv.Atoi(articleStr)
	if err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	articleVersion, err := service.GetLatestArticle(c, int32(articleId))
	if err != nil {
		log.Errorf("获取文章失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "文章不存在", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", articleVersion)
}

// GetArticleListHandler 获取文章列表
func GetArticleListHandler(c *gin.Context) {
	// 直接从URL参数获取
	page := c.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSize := c.DefaultQuery("page_size", "10")
	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}
	if pageSizeNum > 50 {
		pageSizeNum = 50 // 限制最大页面大小
	}

	category := c.Query("category")
	sort := c.DefaultQuery("sort", "created_at_desc")
	search := c.Query("search")
	searchType := c.DefaultQuery("search_type", "fuzzy")
	searchFields := c.DefaultQuery("search_fields", "all")

	articleList, err := service.GetArticleList(c, pageNum, pageSizeNum, category, sort, search, searchType, searchFields)
	if err != nil {
		log.Errorf("获取文章列表失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取文章列表失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
}

// GetHotArticleListHandler 获取热门文章列表
func GetHotArticleListHandler(c *gin.Context) {
	articleList, err := service.GetHotArticleList(c)
	if err != nil {
		log.Errorf("获取热门文章列表失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取热门文章列表失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
}

// GetArticleArchiveHandler 获取文章归档统计
func GetArticleArchiveHandler(c *gin.Context) {
	archiveList, err := service.GetArticleArchive(c)
	if err != nil {
		log.Errorf("获取文章归档失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取文章归档失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", archiveList)
}

func CreateArticleHandler(c *gin.Context) {
	articleReq := &message.ArticleRequest{}
	if err := c.ShouldBindJSON(articleReq); err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	log.Infof("CreateArticle req: %v", articleReq)

	err := service.CreateArticle(c, articleReq)
	if err != nil {
		log.Errorf("文章创建失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "文章创建失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "创建成功", nil)
}

// UpdateArticleHandler 更新文章
func UpdateArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.Atoi(articleStr)
	if err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	articleReq := &message.ArticleRequest{}
	if err := c.ShouldBindJSON(articleReq); err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	err = service.UpdateArticle(c, int32(articleID), articleReq)
	if err != nil {
		log.Errorf("文章更新失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "更新成功", nil)
}

// DeleteArticleHandler 删除文章
func DeleteArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.Atoi(articleStr)
	if err != nil {
		log.Errorf("参数绑定失败: %v", err)
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	if err := service.DeleteArticle(c, int32(articleID)); err != nil {
		log.Errorf("文章删除失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "删除成功", nil)
}

// GetAdminArticleListHandler 获取管理员文章列表
func GetAdminArticleListHandler(c *gin.Context) {
	// 直接从URL参数获取
	page := c.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	// 支持 pageSize 和 page_size 两种格式（向后兼容）
	pageSizeStr := c.Query("page_size")
	if pageSizeStr == "" {
		pageSizeStr = c.Query("pageSize")
	}
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	pageSizeNum, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}
	if pageSizeNum > 50 {
		pageSizeNum = 50 // 限制最大页面大小
	}

	category := c.Query("category")
	sort := c.DefaultQuery("sort", "created_at_desc")
	search := c.Query("search")
	searchType := c.DefaultQuery("search_type", "fuzzy")
	searchFields := c.DefaultQuery("search_fields", "all")
	status := c.Query("status") // 管理员可以按状态筛选

	articleList, err := service.GetAdminArticleList(c, pageNum, pageSizeNum, category, sort, search, searchType, searchFields, status)
	if err != nil {
		log.Errorf("获取管理员文章列表失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取文章列表失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
}

// GetCategoryListHandler 获取分类列表
func GetCategoryListHandler(c *gin.Context) {
	categoryList, err := service.GetCategoryList(c)
	if err != nil {
		log.Errorf("获取分类列表失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取分类列表失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", categoryList)
}

// SearchArticlesHandler 搜索文章
// func SearchArticlesHandler(c *gin.Context) {
// 	keyword := c.Query("keyword")
// 	if keyword == "" {
// 		message.SendMsg(c, http.StatusBadRequest, "搜索关键词不能为空", nil)
// 		return
// 	}

// 	limitStr := c.DefaultQuery("limit", "10")
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	offsetStr := c.DefaultQuery("offset", "0")
// 	offset, err := strconv.Atoi(offsetStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	// 计算页码
// 	page := int32(offset/limit + 1)
// 	pageSize := int32(limit)

// 	articleList, err := service.SearchArticles(c, keyword, page, pageSize)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusInternalServerError, "搜索失败", nil)
// 		return
// 	}

// 	message.SendMsg(c, http.StatusOK, "搜索成功", articleList)
// }

// // GetPopularArticlesHandler 获取热门文章
// func GetPopularArticlesHandler(c *gin.Context) {
// 	limitStr := c.DefaultQuery("limit", "10")
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	articles, err := service.GetPopularArticles(c, int32(limit))
// 	if err != nil {
// 		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
// 		return
// 	}

// 	message.SendMsg(c, http.StatusOK, "获取成功", articles)
// }

// // GetUserArticlesHandler 获取用户的文章列表
// func GetUserArticlesHandler(c *gin.Context) {
// 	userStr := c.Param("user_id")
// 	userID, err := strconv.Atoi(userStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	limitStr := c.DefaultQuery("limit", "10")
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	offsetStr := c.DefaultQuery("offset", "0")
// 	offset, err := strconv.Atoi(offsetStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	// 计算页码
// 	page := int32(offset/limit + 1)
// 	pageSize := int32(limit)

// 	articleList, err := service.GetArticlesByUser(c, int32(userID), page, pageSize)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
// 		return
// 	}

// 	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
// }

// func GetCategoriesHandler(c *gin.Context) {
// 	categories, err := service.GetCategories(c)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
// 		return
// 	}
// 	message.SendMsg(c, http.StatusOK, "获取成功", categories)
// }

// func GetArticlesByCategoryHandler(c *gin.Context) {
// 	categoryStr := c.Param("category_id")
// 	categoryID, err := strconv.Atoi(categoryStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	limitStr := c.DefaultQuery("limit", "10")
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	offsetStr := c.DefaultQuery("offset", "0")
// 	offset, err := strconv.Atoi(offsetStr)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
// 		return
// 	}

// 	// 计算页码
// 	page := int32(offset/limit + 1)
// 	pageSize := int32(limit)

// 	articleList, err := service.GetArticlesByCategory(c, int32(categoryID), page, pageSize)
// 	if err != nil {
// 		message.SendMsg(c, http.StatusBadRequest, err.Error(), nil)
// 		return
// 	}

// 	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
// }
