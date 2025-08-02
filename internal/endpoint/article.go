package endpoint

import (
	"blog/internal/message"
	"blog/internal/service"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleId, err := strconv.Atoi(articleStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	article, err := service.GetArticle(c, int32(articleId))
	if err != nil {
		message.SendMsg(c, http.StatusUnauthorized, "文章不存在", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", article)
}

func GetArticleListHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	// 计算页码
	page := int32(offset/limit + 1)
	pageSize := int32(limit)

	articleList, err := service.GetArticleList(c, page, pageSize, nil, "")
	if err != nil {
		message.SendMsg(c, http.StatusUnauthorized, "文章不存在", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
}

// UploadImageHandler 上传图片
func UploadImageHandler(c *gin.Context) {
	// 检查用户是否登录
	_, exists := c.Get("user_id")
	if !exists {
		message.SendMsg(c, http.StatusUnauthorized, "请先登录", nil)
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "获取文件失败", nil)
		return
	}
	defer file.Close()

	// 检查文件类型
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedTypes[ext] {
		message.SendMsg(c, http.StatusBadRequest, "不支持的文件类型，请上传 jpg、png、gif 或 webp 格式的图片", nil)
		return
	}

	// 检查文件大小（限制为5MB）
	if header.Size > 5*1024*1024 {
		message.SendMsg(c, http.StatusBadRequest, "文件大小不能超过5MB", nil)
		return
	}

	// 生成唯一文件名
	uuidStr := uuid.New().String()
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102_150405"), uuidStr[:8], ext)

	// 创建保存目录
	uploadDir := "internal/endpoint/web/static/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "创建目录失败", nil)
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "创建文件失败", nil)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "保存文件失败", nil)
		return
	}

	// 返回图片访问URL
	imageURL := fmt.Sprintf("/static/images/%s", filename)
	message.SendMsg(c, http.StatusOK, "上传成功", map[string]string{
		"url":      imageURL,
		"filename": filename,
		"original": header.Filename,
	})
}

type ArticleRequest struct {
	Title       string   `json:"title" binding:"required,min=1,max=100"`
	Content     string   `json:"content" binding:"required,min=10"`
	Description string   `json:"description" binding:"max=255"`
	Category    string   `json:"category" binding:"required"`
	Tags        []string `json:"tags"`
}

type UpdateArticleRequest struct {
	Title       string   `json:"title" binding:"omitempty,min=1,max=100"`
	Content     string   `json:"content" binding:"omitempty,min=10"`
	Description string   `json:"description" binding:"omitempty,max=255"`
	Category    string   `json:"category" binding:"omitempty"`
	Tags        []string `json:"tags"`
}

func CreateArticleHandler(c *gin.Context) {
	articleReq := &ArticleRequest{}
	if err := c.ShouldBindJSON(articleReq); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	req := &service.CreateArticleRequest{
		Title:        articleReq.Title,
		Content:      articleReq.Content,
		Description:  articleReq.Description,
		CategoryName: articleReq.Category,
		Tags:         articleReq.Tags,
	}

	article, err := service.CreateArticle(c, req)
	if err != nil {
		message.SendMsg(c, http.StatusUnauthorized, "文章创建失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "创建成功", article)
}

// UpdateArticleHandler 更新文章
func UpdateArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.Atoi(articleStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	updateReq := &UpdateArticleRequest{}
	if err = c.ShouldBindJSON(updateReq); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	// 构建service层的更新请求
	req := &service.UpdateArticleRequest{
		Title:        updateReq.Title,
		Content:      updateReq.Content,
		Description:  updateReq.Description,
		CategoryName: updateReq.Category,
		Tags:         updateReq.Tags,
	}

	article, err := service.UpdateArticle(c, int32(articleID), req)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "更新成功", article)
}

// DeleteArticleHandler 删除文章
func DeleteArticleHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.Atoi(articleStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	if err := service.DeleteArticle(c, int32(articleID)); err != nil {
		message.SendMsg(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "删除成功", nil)
}

// SearchArticlesHandler 搜索文章
func SearchArticlesHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		message.SendMsg(c, http.StatusBadRequest, "搜索关键词不能为空", nil)
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	// 计算页码
	page := int32(offset/limit + 1)
	pageSize := int32(limit)

	articleList, err := service.SearchArticles(c, keyword, page, pageSize)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "搜索失败", nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "搜索成功", articleList)
}

// GetPopularArticlesHandler 获取热门文章
func GetPopularArticlesHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	articles, err := service.GetPopularArticles(c, int32(limit))
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "获取成功", articles)
}

// GetUserArticlesHandler 获取用户的文章列表
func GetUserArticlesHandler(c *gin.Context) {
	userStr := c.Param("user_id")
	userID, err := strconv.Atoi(userStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	// 计算页码
	page := int32(offset/limit + 1)
	pageSize := int32(limit)

	articleList, err := service.GetArticlesByUser(c, int32(userID), page, pageSize)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
}

func GetCategoriesHandler(c *gin.Context) {
	categories, err := service.GetCategories(c)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", categories)
}

func GetArticlesByCategoryHandler(c *gin.Context) {
	categoryStr := c.Param("category_id")
	categoryID, err := strconv.Atoi(categoryStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	// 计算页码
	page := int32(offset/limit + 1)
	pageSize := int32(limit)

	articleList, err := service.GetArticlesByCategory(c, int32(categoryID), page, pageSize)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "获取成功", articleList)
}
