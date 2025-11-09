package service

import (
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiangxu05/logger/v2"
	"gorm.io/gorm"
)

var log = logger.GetLogger()

// ArticleDetail 文章详情响应结构（包含Article和ArticleVersion的完整信息）
type ArticleDetail struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ArticleID   int32  `json:"article_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     int    `json:"version"`
	StoreID     int32  `json:"store_id"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	Views       int    `json:"views"`
}

// GetArticle 获取文章详情
func GetArticle(c *gin.Context, articleId int32, version int) (*ArticleDetail, error) {
	// 增加浏览量
	if err := IncrementArticleViewCount(c, articleId); err != nil {
		// 浏览量更新失败不影响文章获取
		log.Errorf("Failed to increment view count: %v\n", err)
	}

	// 获取文章基础信息
	article, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleId)).
		First()
	if err != nil {
		return nil, err
	}

	// 获取文章版本
	articleVersion, err := dao.ArticleVersion.WithContext(c.Request.Context()).
		Where(dao.ArticleVersion.ArticleID.Eq(articleId), dao.ArticleVersion.Version.Eq(version)).
		First()
	if err != nil {
		return nil, err
	}

	// 组合返回完整的文章详情
	return &ArticleDetail{
		ID:          articleVersion.ID,
		CreatedAt:   articleVersion.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   articleVersion.UpdatedAt.Format("2006-01-02 15:04:05"),
		ArticleID:   article.ID,
		Title:       article.Title,
		Description: articleVersion.Description,
		Version:     articleVersion.Version,
		StoreID:     articleVersion.StoreID,
		Category:    article.Category,
		Tags:        article.Tags,
		Views:       article.Views,
	}, nil
}

// GetLatestArticle 获取最新文章
func GetLatestArticle(c *gin.Context, articleId int32) (*ArticleDetail, error) {
	// 增加浏览量
	if err := IncrementArticleViewCount(c, articleId); err != nil {
		// 浏览量更新失败不影响文章获取
		log.Errorf("Failed to increment view count: %v\n", err)
	}

	// 获取文章基础信息
	article, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleId)).
		First()
	if err != nil {
		return nil, err
	}

	// 获取文章最新版本
	articleVersion, err := dao.ArticleVersion.WithContext(c.Request.Context()).
		Where(dao.ArticleVersion.ArticleID.Eq(articleId), dao.ArticleVersion.Version.Eq(article.Version)).
		First()
	if err != nil {
		return nil, err
	}

	// 组合返回完整的文章详情
	return &ArticleDetail{
		ID:          articleVersion.ID,
		CreatedAt:   articleVersion.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   articleVersion.UpdatedAt.Format("2006-01-02 15:04:05"),
		ArticleID:   article.ID,
		Title:       article.Title,
		Description: articleVersion.Description,
		Version:     articleVersion.Version,
		StoreID:     articleVersion.StoreID,
		Category:    article.Category,
		Tags:        article.Tags,
		Views:       article.Views,
	}, nil
}

// IncrementArticleViewCount 增加文章浏览量（并发安全）
func IncrementArticleViewCount(c *gin.Context, articleId int32) error {
	// 使用原子操作直接增加浏览量，避免并发问题
	_, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleId)).
		Update(dao.Article.Views, gorm.Expr("views + 1"))
	return err
}

// applySearchFilter 应用搜索过滤器
func applySearchFilter(query dao.IArticleDo, search, searchType, searchFields string) dao.IArticleDo {
	if search == "" {
		return query
	}

	// 根据搜索类型构建搜索条件
	switch searchType {
	case "exact":
		// 精确匹配
		return applyExactSearch(query, search, searchFields)
	case "fuzzy":
		// 模糊匹配（默认）
		return applyFuzzySearch(query, search, searchFields)
	default:
		// 默认使用模糊匹配
		return applyFuzzySearch(query, search, searchFields)
	}
}

// applyExactSearch 应用精确搜索
func applyExactSearch(query dao.IArticleDo, search, searchFields string) dao.IArticleDo {
	switch searchFields {
	case "title":
		return query.Where(dao.Q.Article.Title.Eq(search))
	case "content":
		// 暂时在描述字段中搜索
		return query.Where(dao.Q.Article.Description.Eq(search))
	case "tags":
		return query.Where(dao.Q.Article.Tags.Eq(search))
	case "all":
		fallthrough
	default:
		// 在标题、标签、描述中精确搜索
		return query.Where(dao.Q.Article.Title.Eq(search)).
			Or(query.Where(dao.Q.Article.Tags.Eq(search))).
			Or(query.Where(dao.Q.Article.Description.Eq(search)))
	}
}

// applyFuzzySearch 应用模糊搜索
func applyFuzzySearch(query dao.IArticleDo, search, searchFields string) dao.IArticleDo {
	searchPattern := "%" + search + "%"

	switch searchFields {
	case "title":
		return query.Where(dao.Q.Article.Title.Like(searchPattern))
	case "content":
		// 暂时在描述字段中搜索
		return query.Where(dao.Q.Article.Description.Like(searchPattern))
	case "tags":
		return query.Where(dao.Q.Article.Tags.Like(searchPattern))
	case "all":
		fallthrough
	default:
		// 在标题、标签、描述中模糊搜索
		return query.Where(dao.Q.Article.Title.Like(searchPattern)).
			Or(query.Where(dao.Q.Article.Tags.Like(searchPattern))).
			Or(query.Where(dao.Q.Article.Description.Like(searchPattern)))
	}
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context, page, pageSize int, category, sort, search, searchType, searchFields string) (*message.ArticleResponse, error) {
	ctx := c.Request.Context()

	// 构建查询条件
	query := dao.Article.WithContext(ctx)

	// 分类筛选
	if category != "" {
		query = query.Where(dao.Article.Category.Eq(category))
	}

	// 搜索功能
	if search != "" {
		query = applySearchFilter(query, search, searchType, searchFields)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// 排序处理
	switch sort {
	case "created_at_asc":
		query = query.Order(dao.Article.CreatedAt.Asc())
	case "created_at_desc":
		query = query.Order(dao.Article.CreatedAt.Desc())
	case "title_asc":
		query = query.Order(dao.Article.Title.Asc())
	case "title_desc":
		query = query.Order(dao.Article.Title.Desc())
	case "views_desc":
		query = query.Order(dao.Article.Views.Desc())
	case "views_asc":
		query = query.Order(dao.Article.Views.Asc())
	default:
		// 默认按创建时间倒序
		query = query.Order(dao.Article.CreatedAt.Desc())
	}

	// 获取文章列表
	offset := (page - 1) * pageSize
	articles, err := query.
		Limit(pageSize).
		Offset(offset).
		Find()
	if err != nil {
		return nil, err
	}

	// 创建一个新的结构体切片
	var articleStructs []model_def.Article
	for _, article := range articles {
		articleStructs = append(articleStructs, *article)
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	resp := &message.ArticleResponse{
		Articles: articleStructs,
		Pagination: message.PaginationInfo{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  int(total),
		},
	}

	return resp, nil
}

// GetHotArticleList 获取热门文章列表
func GetHotArticleList(c *gin.Context) (*message.ArticleResponse, error) {
	// 获取热门文章列表
	articles, err := dao.Article.WithContext(c.Request.Context()).
		Order(dao.Article.Views.Desc()).
		Limit(3).
		Find()
	if err != nil {
		return nil, err
	}
	// 创建一个新的结构体切片
	var articleStructs []model_def.Article
	for _, article := range articles {
		articleStructs = append(articleStructs, *article)
	}
	resp := &message.ArticleResponse{
		Articles: articleStructs,
	}
	return resp, nil
}

// CreateArticle 创建文章
func CreateArticle(c *gin.Context, req *message.ArticleRequest) error {
	// 获取当前用户ID
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		return errors.New("用户未登录")
	}
	userID := userIdInterface.(int32)

	// 使用事务创建文章
	var article *model_def.Article
	err := dao.Q.Transaction(func(tx *dao.Query) error {
		// 处理分类
		category, err := processArticleCategory(c.Request.Context(), req.Category)
		if err != nil {
			return err
		}
		// 创建文章基本信息
		article = &model_def.Article{
			UserID:      userID,
			Title:       req.Title,
			Description: req.Description,
			Version:     1,
			Category:    category.Name,
			Tags:        req.Tags,
			Views:       0,
		}
		if err := tx.Article.WithContext(c.Request.Context()).Create(article); err != nil {
			return err
		}
		// 创建文章版本
		articleVersion := &model_def.ArticleVersion{
			ArticleID:   article.ID,
			Description: req.Description,
			Version:     1,
			StoreID:     req.StoreID,
		}
		if err := tx.ArticleVersion.WithContext(c.Request.Context()).Create(articleVersion); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func processArticleCategory(ctx context.Context, categoryName string) (*model_def.Category, error) {
	// 去除前后空格
	categoryName = strings.TrimSpace(categoryName)
	if categoryName == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 查找分类
	category, err := dao.Category.WithContext(ctx).
		Where(dao.Category.Name.Eq(categoryName)).
		First()

	// 如果找到分类，直接返回
	if err == nil {
		category.Number++
		if err := dao.Category.WithContext(ctx).Save(category); err != nil {
			return nil, fmt.Errorf("更新分类失败: %v", err)
		}
		return category, nil
	}

	// 如果错误不是"记录未找到"，返回错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询分类失败: %v", err)
	}

	// 分类不存在，创建新分类
	newCategory := &model_def.Category{
		Name:   categoryName,
		Number: 1,
	}

	// 创建分类
	if err = dao.Category.WithContext(ctx).Create(newCategory); err != nil {
		return nil, fmt.Errorf("创建分类失败: %v", err)
	}

	return newCategory, nil
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context, articleID int32, req *message.ArticleRequest) error {
	// 获取当前用户ID
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		return errors.New("用户未登录")
	}
	userID := userIdInterface.(int32)

	// 检查文章是否存在且属于当前用户
	existingArticle, err := dao.Article.WithContext(c.Request.Context()).
		Where(dao.Article.ID.Eq(articleID)).
		Where(dao.Article.UserID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在或无权限修改")
		}
		return err
	}
	// 使用事务更新文章
	err = dao.Q.Transaction(func(tx *dao.Query) error {
		// 更新文章基本信息
		existingArticle.Title = req.Title
		existingArticle.Category = req.Category
		existingArticle.Tags = req.Tags
		existingArticle.Version++
		// 更新文章
		if err := tx.Article.WithContext(c.Request.Context()).Save(existingArticle); err != nil {
			return err
		}
		// 创建文章版本
		articleVersion := &model_def.ArticleVersion{
			ArticleID:   existingArticle.ID,
			Description: req.Description,
			Version:     existingArticle.Version,
			StoreID:     req.StoreID,
		}
		if err = tx.ArticleVersion.WithContext(c.Request.Context()).Create(articleVersion); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
	existingArticle, err := dao.Article.WithContext(c.Request.Context()).
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
		// 删除文章
		_, err = tx.Article.WithContext(c.Request.Context()).
			Where(tx.Article.ID.Eq(articleID)).
			Delete()
		if err != nil {
			return err
		}
		// 删除文章版本
		_, err = tx.ArticleVersion.WithContext(c.Request.Context()).
			Where(tx.ArticleVersion.ArticleID.Eq(articleID)).
			Delete()
		if err != nil {
			return err
		}
		// 检查是否需要删除分类
		category, err := tx.Category.WithContext(c.Request.Context()).
			Where(tx.Category.Name.Eq(existingArticle.Category)).
			First()
		if err != nil {
			return err
		}
		if category.Number == 1 {
			_, err = tx.Category.WithContext(c.Request.Context()).
				Where(tx.Category.ID.Eq(category.ID)).
				Delete()
			if err != nil {
				return err
			}
		} else {
			category.Number--
			if err := tx.Category.WithContext(c.Request.Context()).Save(category); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetAdminArticleList 获取管理员文章列表（包含所有状态的文章）
func GetAdminArticleList(c *gin.Context, page, pageSize int, category, sort, search, searchType, searchFields, status string) (*message.ArticleResponse, error) {
	ctx := c.Request.Context()

	// 构建查询条件
	query := dao.Article.WithContext(ctx)

	// 分类筛选
	if category != "" {
		query = query.Where(dao.Article.Category.Eq(category))
	}

	// 状态筛选（管理员专用）
	if status != "" {
		// 这里可以根据实际需求添加状态字段筛选
		// 目前模型中没有status字段，可以根据其他条件筛选
		// 例如：根据是否有内容、创建时间等判断状态
	}

	// 搜索功能
	if search != "" {
		query = applySearchFilter(query, search, searchType, searchFields)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// 排序处理
	switch sort {
	case "created_at_asc":
		query = query.Order(dao.Article.CreatedAt.Asc())
	case "created_at_desc":
		query = query.Order(dao.Article.CreatedAt.Desc())
	case "updated_at_desc":
		query = query.Order(dao.Article.UpdatedAt.Desc())
	case "updated_at_asc":
		query = query.Order(dao.Article.UpdatedAt.Asc())
	case "title_asc":
		query = query.Order(dao.Article.Title.Asc())
	case "title_desc":
		query = query.Order(dao.Article.Title.Desc())
	case "views_desc":
		query = query.Order(dao.Article.Views.Desc())
	case "views_asc":
		query = query.Order(dao.Article.Views.Asc())
	default:
		// 默认按创建时间倒序
		query = query.Order(dao.Article.CreatedAt.Desc())
	}

	// 获取文章列表
	offset := (page - 1) * pageSize
	articles, err := query.
		Limit(pageSize).
		Offset(offset).
		Find()
	if err != nil {
		return nil, err
	}

	// 创建一个新的结构体切片
	var articleStructs []model_def.Article
	for _, article := range articles {
		articleStructs = append(articleStructs, *article)
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	resp := &message.ArticleResponse{
		Articles: articleStructs,
		Pagination: message.PaginationInfo{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  int(total),
		},
	}

	return resp, nil
}

func GetCategoryList(c *gin.Context) ([]*model_def.Category, error) {
	categoryList, err := dao.Category.WithContext(c.Request.Context()).Find()
	if err != nil {
		return nil, err
	}
	return categoryList, nil
}

// // GetArticlesByUser 获取用户的文章列表
// func GetArticlesByUser(c *gin.Context, userID int32, page, pageSize int32) (*ArticleListResponse, error) {
// 	ctx := c.Request.Context()
// 	offset := (page - 1) * pageSize

// 	// 获取总数
// 	total, err := dao.Article.WithContext(ctx).
// 		Where(dao.Article.UserID.Eq(userID)).
// 		Count()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 获取文章列表
// 	articles, err := dao.Article.WithContext(ctx).
// 		Preload(dao.Article.User).
// 		Preload(dao.Article.Category).
// 		Preload(dao.Article.Tags).
// 		Where(dao.Article.UserID.Eq(userID)).
// 		Order(dao.Article.CreatedAt.Desc()).
// 		Limit(int(pageSize)).
// 		Offset(int(offset)).
// 		Find()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ArticleListResponse{
// 		Articles: articles,
// 		Total:    total,
// 		Page:     page,
// 		PageSize: pageSize,
// 	}, nil
// }

// // SearchArticles 搜索文章
// func SearchArticles(c *gin.Context, keyword string, page, pageSize int32) (*ArticleListResponse, error) {
// 	ctx := c.Request.Context()
// 	offset := (page - 1) * pageSize
// 	keyword = "%" + strings.TrimSpace(keyword) + "%"

// 	// 构建搜索查询
// 	query := dao.Article.WithContext(ctx).
// 		Preload(dao.Article.User).
// 		Preload(dao.Article.Category).
// 		Preload(dao.Article.Tags).
// 		Where(dao.Article.Title.Like(keyword))

// 	// 获取总数
// 	total, err := query.Count()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 获取搜索结果
// 	articles, err := query.
// 		Order(dao.Article.CreatedAt.Desc()).
// 		Limit(int(pageSize)).
// 		Offset(int(offset)).
// 		Find()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ArticleListResponse{
// 		Articles: articles,
// 		Total:    total,
// 		Page:     page,
// 		PageSize: pageSize,
// 	}, nil
// }

// // GetPopularArticles 获取热门文章
// func GetPopularArticles(c *gin.Context, limit int32) ([]*model_def.Article, error) {
// 	articles, err := dao.Article.WithContext(c.Request.Context()).
// 		Preload(dao.Article.User).
// 		Preload(dao.Article.Category).
// 		Preload(dao.Article.Tags).
// 		Order(dao.Article.Views.Desc()).
// 		Limit(int(limit)).
// 		Find()
// 	return articles, err
// }

// // GetCategories 获取所有分类
// func GetCategories(c *gin.Context) ([]*model_def.Category, error) {
// 	categories, err := dao.Category.WithContext(c.Request.Context()).
// 		Order(dao.Category.Name).
// 		Find()
// 	return categories, err
// }

// // GetArticlesByCategory 根据分类ID获取文章列表
// func GetArticlesByCategory(c *gin.Context, categoryID int32, page, pageSize int32) (*ArticleListResponse, error) {
// 	ctx := c.Request.Context()
// 	offset := (page - 1) * pageSize

// 	// 验证分类是否存在
// 	_, err := dao.Category.WithContext(ctx).
// 		Where(dao.Category.ID.Eq(categoryID)).
// 		First()
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("分类不存在")
// 		}
// 		return nil, err
// 	}

// 	// 获取总数
// 	total, err := dao.Article.WithContext(ctx).
// 		Where(dao.Article.CategoryID.Eq(categoryID)).
// 		Count()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 获取文章列表
// 	articles, err := dao.Article.WithContext(ctx).
// 		Preload(dao.Article.User).
// 		Preload(dao.Article.Category).
// 		Preload(dao.Article.Tags).
// 		Where(dao.Article.CategoryID.Eq(categoryID)).
// 		Order(dao.Article.CreatedAt.Desc()).
// 		Limit(int(pageSize)).
// 		Offset(int(offset)).
// 		Find()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ArticleListResponse{
// 		Articles: articles,
// 		Total:    total,
// 		Page:     page,
// 		PageSize: pageSize,
// 	}, nil
// }

// // 辅助函数：处理文章分类
// func processArticleCategory(ctx context.Context, categoryName string) (*model_def.Category, error) {
// 	categoryName = strings.TrimSpace(categoryName)
// 	if categoryName == "" {
// 		return nil, errors.New("分类名称不能为空")
// 	}

// 	// 查找分类
// 	category, err := dao.Category.WithContext(ctx).
// 		Where(dao.Category.Name.Eq(categoryName)).
// 		First()
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// 分类不存在，创建新分类
// 			newCategory := &model_def.Category{
// 				Name: categoryName,
// 			}
// 			if err = dao.Category.WithContext(ctx).Create(newCategory); err != nil {
// 				return nil, fmt.Errorf("创建分类失败: %v", err)
// 			}
// 			return newCategory, nil
// 		} else {
// 			return nil, fmt.Errorf("查询分类失败: %v", err)
// 		}
// 	}

// 	return category, nil
// }

// // 辅助函数：处理文章标签
// func processArticleTags(ctx context.Context, tagNames []string) ([]*model_def.Tag, error) {
// 	if len(tagNames) == 0 {
// 		return nil, nil
// 	}

// 	var tags []*model_def.Tag
// 	for _, tagName := range tagNames {
// 		tagName = strings.TrimSpace(tagName)
// 		if tagName == "" {
// 			continue
// 		}

// 		// 查找或创建标签
// 		tag, err := dao.Tag.WithContext(ctx).
// 			Where(dao.Tag.Name.Eq(tagName)).
// 			First()
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				// 标签不存在，创建新标签
// 				newTag := &model_def.Tag{Name: tagName}
// 				if err := dao.Tag.WithContext(ctx).Create(newTag); err != nil {
// 					return nil, err
// 				}
// 				tags = append(tags, newTag)
// 			} else {
// 				return nil, err
// 			}
// 		} else {
// 			tags = append(tags, tag)
// 		}
// 	}

// 	return tags, nil
// }

// // 辅助函数：获取文章但不增加浏览量
// func GetArticleWithoutViewIncrement(c *gin.Context, articleID int32) (*model_def.Article, error) {
// 	// 获取文章基本信息
// 	article, err := dao.Article.WithContext(c.Request.Context()).
// 		Preload(dao.Article.User).
// 		Preload(dao.Article.Category).
// 		Preload(dao.Article.Tags).
// 		Where(dao.Article.ID.Eq(articleID)).
// 		First()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 获取文章内容
// 	articleContent, err := dao.ArticleContent.WithContext(c.Request.Context()).
// 		Where(dao.ArticleContent.ArticleID.Eq(articleID)).
// 		First()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 将内容设置到文章对象中
// 	article.Content = *articleContent

// 	return article, nil
// }
