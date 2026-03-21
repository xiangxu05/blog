package service

import (
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetArticleSocial 文章点赞数与当前用户状态
func GetArticleSocial(c *gin.Context, articleID int32, userID *int32) (*message.ArticleSocial, error) {
	ctx := c.Request.Context()

	_, err := dao.Article.WithContext(ctx).Where(dao.Article.ID.Eq(articleID)).First()
	if err != nil {
		return nil, err
	}

	n, err := dao.ArticleLike.WithContext(ctx).Where(dao.ArticleLike.ArticleID.Eq(articleID)).Count()
	if err != nil {
		return nil, err
	}
	out := &message.ArticleSocial{LikeCount: int(n)}
	if userID != nil {
		ln, err := dao.ArticleLike.WithContext(ctx).Where(
			dao.ArticleLike.ArticleID.Eq(articleID),
			dao.ArticleLike.UserID.Eq(*userID),
		).Count()
		if err != nil {
			return nil, err
		}
		out.Liked = ln > 0
		bn, err := dao.ArticleBookmark.WithContext(ctx).Where(
			dao.ArticleBookmark.ArticleID.Eq(articleID),
			dao.ArticleBookmark.UserID.Eq(*userID),
		).Count()
		if err != nil {
			return nil, err
		}
		out.Bookmarked = bn > 0
	}
	return out, nil
}

// LikeArticle 点赞（gorm.Model 含软删：取消赞若为软删则再次 Create 会撞唯一索引，需恢复或硬删）
func LikeArticle(c *gin.Context, articleID int32, userID int32) error {
	ctx := c.Request.Context()
	if _, err := dao.Article.WithContext(ctx).Where(dao.Article.ID.Eq(articleID)).First(); err != nil {
		return err
	}
	row, err := dao.ArticleLike.WithContext(ctx).Unscoped().Where(
		dao.ArticleLike.ArticleID.Eq(articleID),
		dao.ArticleLike.UserID.Eq(userID),
	).First()
	if err == nil {
		if row.DeletedAt.Valid {
			_, err := dao.ArticleLike.WithContext(ctx).Unscoped().Where(
				dao.ArticleLike.ID.Eq(row.ID),
			).Updates(map[string]interface{}{"deleted_at": nil})
			return err
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return dao.ArticleLike.WithContext(ctx).Create(&model_def.ArticleLike{
		UserID:    userID,
		ArticleID: articleID,
	})
}

// UnlikeArticle 取消赞（物理删除，避免与唯一索引 + 软删组合导致无法再次点赞）
func UnlikeArticle(c *gin.Context, articleID int32, userID int32) error {
	ctx := c.Request.Context()
	_, err := dao.ArticleLike.WithContext(ctx).Unscoped().Where(
		dao.ArticleLike.ArticleID.Eq(articleID),
		dao.ArticleLike.UserID.Eq(userID),
	).Delete()
	return err
}

// BookmarkArticle 收藏（逻辑同文章赞，避免软删 + 唯一索引）
func BookmarkArticle(c *gin.Context, articleID int32, userID int32) error {
	ctx := c.Request.Context()
	if _, err := dao.Article.WithContext(ctx).Where(dao.Article.ID.Eq(articleID)).First(); err != nil {
		return err
	}
	row, err := dao.ArticleBookmark.WithContext(ctx).Unscoped().Where(
		dao.ArticleBookmark.ArticleID.Eq(articleID),
		dao.ArticleBookmark.UserID.Eq(userID),
	).First()
	if err == nil {
		if row.DeletedAt.Valid {
			_, err := dao.ArticleBookmark.WithContext(ctx).Unscoped().Where(
				dao.ArticleBookmark.ID.Eq(row.ID),
			).Updates(map[string]interface{}{"deleted_at": nil})
			return err
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return dao.ArticleBookmark.WithContext(ctx).Create(&model_def.ArticleBookmark{
		UserID:    userID,
		ArticleID: articleID,
	})
}

// UnbookmarkArticle 取消收藏
func UnbookmarkArticle(c *gin.Context, articleID int32, userID int32) error {
	ctx := c.Request.Context()
	_, err := dao.ArticleBookmark.WithContext(ctx).Unscoped().Where(
		dao.ArticleBookmark.ArticleID.Eq(articleID),
		dao.ArticleBookmark.UserID.Eq(userID),
	).Delete()
	return err
}

// UpdateComment 编辑自己的评论
func UpdateComment(c *gin.Context, commentID uint, userID int32, content string) error {
	ctx := c.Request.Context()
	content = strings.TrimSpace(content)
	if len([]rune(content)) < 1 || len([]rune(content)) > 2000 {
		return errors.New("评论内容长度无效")
	}
	row, err := dao.Comment.WithContext(ctx).Where(dao.Comment.ID.Eq(commentID)).First()
	if err != nil {
		return err
	}
	if row.UserID != userID {
		return errors.New("forbidden")
	}
	row.Content = content
	return dao.Comment.WithContext(ctx).Save(row)
}

// ToggleCommentLike 点赞评论（幂等；处理软删残留）
func ToggleCommentLike(c *gin.Context, commentID uint, userID int32) error {
	ctx := c.Request.Context()
	if _, err := dao.Comment.WithContext(ctx).Where(dao.Comment.ID.Eq(commentID)).First(); err != nil {
		return err
	}
	row, err := dao.CommentLike.WithContext(ctx).Unscoped().Where(
		dao.CommentLike.CommentID.Eq(commentID),
		dao.CommentLike.UserID.Eq(userID),
	).First()
	if err == nil {
		if row.DeletedAt.Valid {
			_, err := dao.CommentLike.WithContext(ctx).Unscoped().Where(
				dao.CommentLike.ID.Eq(row.ID),
			).Updates(map[string]interface{}{"deleted_at": nil})
			return err
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return dao.CommentLike.WithContext(ctx).Create(&model_def.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	})
}

// RemoveCommentLike 取消评论赞
func RemoveCommentLike(c *gin.Context, commentID uint, userID int32) error {
	ctx := c.Request.Context()
	_, err := dao.CommentLike.WithContext(ctx).Unscoped().Where(
		dao.CommentLike.CommentID.Eq(commentID),
		dao.CommentLike.UserID.Eq(userID),
	).Delete()
	return err
}

// CommentLikeCount 评论点赞数
func CommentLikeCount(c *gin.Context, commentID uint) (int64, error) {
	return dao.CommentLike.WithContext(c.Request.Context()).Where(dao.CommentLike.CommentID.Eq(commentID)).Count()
}

// UserLikedComment 当前用户是否已赞
func UserLikedComment(c *gin.Context, commentID uint, userID int32) (bool, error) {
	n, err := dao.CommentLike.WithContext(c.Request.Context()).Where(
		dao.CommentLike.CommentID.Eq(commentID),
		dao.CommentLike.UserID.Eq(userID),
	).Count()
	return n > 0, err
}

// ListMyBookmarks 当前用户收藏的文章（按收藏时间倒序）
func ListMyBookmarks(c *gin.Context, userID int32, page, pageSize int) (*message.BookmarkListResponse, error) {
	ctx := c.Request.Context()
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	q := dao.ArticleBookmark.WithContext(ctx).Where(dao.ArticleBookmark.UserID.Eq(userID))
	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	bookmarks, err := dao.ArticleBookmark.WithContext(ctx).
		Where(dao.ArticleBookmark.UserID.Eq(userID)).
		Order(dao.ArticleBookmark.CreatedAt.Desc()).
		Offset(offset).
		Limit(pageSize).
		Find()
	if err != nil {
		return nil, err
	}

	items := make([]message.BookmarkArticleItem, 0, len(bookmarks))
	for _, bm := range bookmarks {
		a, err := dao.Article.WithContext(ctx).Where(dao.Article.ID.Eq(bm.ArticleID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, err
		}
		items = append(items, message.BookmarkArticleItem{
			Article:      *a,
			BookmarkedAt: bm.CreatedAt,
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &message.BookmarkListResponse{
		Articles: items,
		Pagination: message.PaginationInfo{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  int(total),
		},
	}, nil
}
