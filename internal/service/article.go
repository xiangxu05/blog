package service

import (
	"blog/dao"
	"blog/model_def"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateArticleRequest 创建文章请求结构
type CreateArticleRequest struct {
	Title        string   `json:"title" binding:"required,min=1,max=100"`
	Content      string   `json:"content" binding:"required,min=10"`
	Description  string   `json:"description" binding:"max=255"`
	CategoryName string   `json:"category_name" binding:"required,min=1,max=32"`
	Tags         []string `json:"tags"`
}

// UpdateArticleRequest 更新文章请求结构
type UpdateArticleRequest struct {
	Title        string   `json:"title" binding:"omitempty,min=1,max=100"`
	Content      string   `json:"content" binding:"omitempty,min=10"`
	Description  string   `json:"description" binding:"omitempty,max=255"`
	CategoryName string   `json:"category_name" binding:"omitempty,min=1,max=32"`
	Tags         []string `json:"tags"`
}

// ArticleListResponse 文章列表响应结构
type ArticleListResponse struct {
	Articles []*model_def.Article `json:"articles"`
	Total    int64                `json:"total"`
	Page     int32                `json:"page"`
	PageSize int32                `json:"page_size"`
}

// GetArticle 获取文章详情
func GetArticle(c *gin.Context, articleId int32) (*model_def.Article, error) {
	// 获取文章基本信息
	article, err := dao.Article.WithContext(c.Request.Context()).
		Preload(dao.Article.User).
		Preload(dao.Article.Content).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags).
		Where(dao.Article.ID.Eq(articleId)).
		First()
	if err != nil {
		return nil, err
	}

	// 增加浏览量
	if err := IncrementArticleViewCount(c, articleId); err != nil {
		// 浏览量更新失败不影响文章获取
		fmt.Printf("Failed to increment view count: %v\n", err)
	}

	return article, nil
}

// IncrementArticleViewCount 增加文章浏览量（并发安全）
func IncrementArticleViewCount(c *gin.Context, articleId int32) error {
	// 使用原子操作直接增加浏览量，避免并发问题
	_, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleId)).
		Update(dao.Article.Views, gorm.Expr("views + 1"))
	return err
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context, page, pageSize int32, categoryID *int32, tag string) (*ArticleListResponse, error) {
	ctx := c.Request.Context()
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := dao.Article.WithContext(ctx).
		Preload(dao.Article.User).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags)

	// 按分类筛选
	if categoryID != nil && *categoryID > 0 {
		query = query.Where(dao.Article.CategoryID.Eq(*categoryID))
	}

	// 按标签筛选
	if tag != "" {
		// 先查找标签ID
		tagRecord, err := dao.Tag.WithContext(ctx).Where(dao.Tag.Name.Eq(tag)).First()
		if err != nil {
			// 如果标签不存在，返回空结果
			return &ArticleListResponse{
				Articles: []*model_def.Article{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
			}, nil
		}

		// 查找包含该标签的文章ID
		articleTagRecords, err := dao.ArticleTag.WithContext(ctx).
			Where(dao.ArticleTag.TagID.Eq(int32(tagRecord.ID))).
			Find()
		if err != nil {
			return nil, err
		}

		// 提取文章ID
		var articleIDs []int32
		for _, at := range articleTagRecords {
			articleIDs = append(articleIDs, at.ArticleID)
		}

		if len(articleIDs) == 0 {
			return &ArticleListResponse{
				Articles: []*model_def.Article{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
			}, nil
		}

		// 筛选文章
		query = query.Where(dao.Article.ID.In(articleIDs...))
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// 获取文章列表
	articles, err := query.
		Order(dao.Article.CreatedAt.Desc()).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		return nil, err
	}

	return &ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateArticle 创建文章
func CreateArticle(c *gin.Context, req *CreateArticleRequest) (*model_def.Article, error) {
	// 获取当前用户ID
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("用户未登录")
	}
	userID := userIdInterface.(int32)

	// 查找或创建分类
	category, err := processArticleCategory(c.Request.Context(), req.CategoryName)
	if err != nil {
		return nil, err
	}

	// 处理标签
	tags, err := processArticleTags(c.Request.Context(), req.Tags)
	if err != nil {
		return nil, err
	}

	// 使用事务创建文章
	var article *model_def.Article
	err = dao.Q.Transaction(func(tx *dao.Query) error {
		// 创建文章基本信息
		article = &model_def.Article{
			UserID:     userID,
			Title:      req.Title,
			CategoryID: int32(category.ID),
			Views:      0,
			Version:    1,
		}

		if err := tx.Article.WithContext(c.Request.Context()).Create(article); err != nil {
			return err
		}

		// 创建文章内容
		articleContent := &model_def.ArticleContent{
			ArticleID:   article.ID,
			Content:     req.Content,
			Description: req.Description,
			Version:     1,
		}

		if err := tx.ArticleContent.WithContext(c.Request.Context()).Create(articleContent); err != nil {
			return err
		}

		// 关联标签
		if len(tags) > 0 {
			for _, tag := range tags {
				articleTag := &model_def.ArticleTag{
					ArticleID: article.ID,
					TagID:     int32(tag.ID),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := tx.ArticleTag.WithContext(c.Request.Context()).Create(articleTag); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 重新获取完整的文章信息
	return GetArticleWithoutViewIncrement(c, article.ID)
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context, articleID int32, req *UpdateArticleRequest) (*model_def.Article, error) {
	// 获取当前用户ID
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("用户未登录")
	}
	userID := userIdInterface.(int32)

	// 检查文章是否存在且属于当前用户
	existingArticle, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleID)).
		Where(dao.Article.UserID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在或无权限修改")
		}
		return nil, err
	}

	// 处理分类（如果提供）
	var category *model_def.Category
	if req.CategoryName != "" {
		category, err = processArticleCategory(c.Request.Context(), req.CategoryName)
		if err != nil {
			return nil, err
		}
	}

	// 处理标签
	var tags []*model_def.Tag
	if req.Tags != nil {
		tags, err = processArticleTags(c.Request.Context(), req.Tags)
		if err != nil {
			return nil, err
		}
	}

	// 使用事务更新文章
	err = dao.Q.Transaction(func(tx *dao.Query) error {
		// 更新文章基本信息
		updateData := make(map[string]interface{})
		if req.Title != "" {
			updateData["title"] = req.Title
		}
		if category != nil {
			updateData["category_id"] = int32(category.ID)
		}
		updateData["version"] = existingArticle.Version + 1
		updateData["updated_at"] = time.Now()

		if len(updateData) > 0 {
			_, err := tx.Article.WithContext(c.Request.Context()).
				Where(tx.Article.ID.Eq(articleID)).
				Updates(updateData)
			if err != nil {
				return err
			}
		}

		// 更新文章内容
		if req.Content != "" || req.Description != "" {
			contentUpdateData := make(map[string]interface{})
			if req.Content != "" {
				contentUpdateData["content"] = req.Content
			}
			if req.Description != "" {
				contentUpdateData["description"] = req.Description
			}
			contentUpdateData["version"] = existingArticle.Version + 1
			contentUpdateData["updated_at"] = time.Now()

			_, err := tx.ArticleContent.WithContext(c.Request.Context()).
				Where(tx.ArticleContent.ArticleID.Eq(articleID)).
				Updates(contentUpdateData)
			if err != nil {
				return err
			}
		}

		// 更新标签关联
		if req.Tags != nil {
			// 删除现有标签关联
			_, err := tx.ArticleTag.WithContext(c.Request.Context()).
				Where(tx.ArticleTag.ArticleID.Eq(articleID)).
				Delete()
			if err != nil {
				return err
			}

			// 创建新的标签关联
			for _, tag := range tags {
				articleTag := &model_def.ArticleTag{
					ArticleID: articleID,
					TagID:     int32(tag.ID),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := tx.ArticleTag.WithContext(c.Request.Context()).Create(articleTag); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 返回更新后的文章
	return GetArticleWithoutViewIncrement(c, articleID)
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context, articleID int32) error {
	// 获取当前用户ID
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		return errors.New("用户未登录")
	}
	userID := userIdInterface.(int32)

	// 检查文章是否存在且属于当前用户
	_, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleID)).
		Where(dao.Article.UserID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在或无权限删除")
		}
		return err
	}

	// 使用事务删除文章及相关数据
	return dao.Q.Transaction(func(tx *dao.Query) error {
		// 删除标签关联
		_, err := tx.ArticleTag.WithContext(c.Request.Context()).
			Where(tx.ArticleTag.ArticleID.Eq(articleID)).
			Delete()
		if err != nil {
			return err
		}

		// 删除文章内容
		_, err = tx.ArticleContent.WithContext(c.Request.Context()).
			Where(tx.ArticleContent.ArticleID.Eq(articleID)).
			Delete()
		if err != nil {
			return err
		}

		// 删除文章
		_, err = tx.Article.WithContext(c.Request.Context()).
			Where(tx.Article.ID.Eq(articleID)).
			Delete()
		return err
	})
}

// GetArticlesByUser 获取用户的文章列表
func GetArticlesByUser(c *gin.Context, userID int32, page, pageSize int32) (*ArticleListResponse, error) {
	ctx := c.Request.Context()
	offset := (page - 1) * pageSize

	// 获取总数
	total, err := dao.Article.WithContext(ctx).
		Where(dao.Article.UserID.Eq(userID)).
		Count()
	if err != nil {
		return nil, err
	}

	// 获取文章列表
	articles, err := dao.Article.WithContext(ctx).
		Preload(dao.Article.User).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags).
		Where(dao.Article.UserID.Eq(userID)).
		Order(dao.Article.CreatedAt.Desc()).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		return nil, err
	}

	return &ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// SearchArticles 搜索文章
func SearchArticles(c *gin.Context, keyword string, page, pageSize int32) (*ArticleListResponse, error) {
	ctx := c.Request.Context()
	offset := (page - 1) * pageSize
	keyword = "%" + strings.TrimSpace(keyword) + "%"

	// 构建搜索查询
	query := dao.Article.WithContext(ctx).
		Preload(dao.Article.User).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags).
		Where(dao.Article.Title.Like(keyword))

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// 获取搜索结果
	articles, err := query.
		Order(dao.Article.CreatedAt.Desc()).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		return nil, err
	}

	return &ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetPopularArticles 获取热门文章
func GetPopularArticles(c *gin.Context, limit int32) ([]*model_def.Article, error) {
	articles, err := dao.Article.WithContext(c.Request.Context()).
		Preload(dao.Article.User).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags).
		Order(dao.Article.Views.Desc()).
		Limit(int(limit)).
		Find()
	return articles, err
}

// GetCategories 获取所有分类
func GetCategories(c *gin.Context) ([]*model_def.Category, error) {
	categories, err := dao.Category.WithContext(c.Request.Context()).
		Order(dao.Category.Name).
		Find()
	return categories, err
}

// GetArticlesByCategory 根据分类ID获取文章列表
func GetArticlesByCategory(c *gin.Context, categoryID int32, page, pageSize int32) (*ArticleListResponse, error) {
	ctx := c.Request.Context()
	offset := (page - 1) * pageSize

	// 验证分类是否存在
	_, err := dao.Category.WithContext(ctx).
		Where(dao.Category.ID.Eq(categoryID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	// 获取总数
	total, err := dao.Article.WithContext(ctx).
		Where(dao.Article.CategoryID.Eq(categoryID)).
		Count()
	if err != nil {
		return nil, err
	}

	// 获取文章列表
	articles, err := dao.Article.WithContext(ctx).
		Preload(dao.Article.User).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags).
		Where(dao.Article.CategoryID.Eq(categoryID)).
		Order(dao.Article.CreatedAt.Desc()).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		return nil, err
	}

	return &ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// 辅助函数：处理文章分类
func processArticleCategory(ctx context.Context, categoryName string) (*model_def.Category, error) {
	categoryName = strings.TrimSpace(categoryName)
	if categoryName == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 查找分类
	category, err := dao.Category.WithContext(ctx).
		Where(dao.Category.Name.Eq(categoryName)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 分类不存在，创建新分类
			newCategory := &model_def.Category{
				Name: categoryName,
			}
			if err = dao.Category.WithContext(ctx).Create(newCategory); err != nil {
				return nil, fmt.Errorf("创建分类失败: %v", err)
			}
			return newCategory, nil
		} else {
			return nil, fmt.Errorf("查询分类失败: %v", err)
		}
	}

	return category, nil
}

// 辅助函数：处理文章标签
func processArticleTags(ctx context.Context, tagNames []string) ([]*model_def.Tag, error) {
	if len(tagNames) == 0 {
		return nil, nil
	}

	var tags []*model_def.Tag
	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// 查找或创建标签
		tag, err := dao.Tag.WithContext(ctx).
			Where(dao.Tag.Name.Eq(tagName)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 标签不存在，创建新标签
				newTag := &model_def.Tag{Name: tagName}
				if err := dao.Tag.WithContext(ctx).Create(newTag); err != nil {
					return nil, err
				}
				tags = append(tags, newTag)
			} else {
				return nil, err
			}
		} else {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// 辅助函数：获取文章但不增加浏览量
func GetArticleWithoutViewIncrement(c *gin.Context, articleID int32) (*model_def.Article, error) {
	// 获取文章基本信息
	article, err := dao.Article.WithContext(c.Request.Context()).
		Preload(dao.Article.User).
		Preload(dao.Article.Category).
		Preload(dao.Article.Tags).
		Where(dao.Article.ID.Eq(articleID)).
		First()
	if err != nil {
		return nil, err
	}

	// 获取文章内容
	articleContent, err := dao.ArticleContent.WithContext(c.Request.Context()).
		Where(dao.ArticleContent.ArticleID.Eq(articleID)).
		First()
	if err != nil {
		return nil, err
	}

	// 将内容设置到文章对象中
	article.Content = *articleContent

	return article, nil
}
