package message

import (
	"blog/model_def"
	"time"
)

type ArticleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	StoreID     int32  `json:"store_id"`
}

type ArticleResponse struct {
	Articles   []model_def.Article `json:"articles"`
	Pagination PaginationInfo      `json:"pagination"`
}

// BookmarkArticleItem 收藏列表单项（文章字段 + 收藏时间）
type BookmarkArticleItem struct {
	model_def.Article
	BookmarkedAt time.Time `json:"bookmarked_at"`
}

// BookmarkListResponse 我的收藏列表
type BookmarkListResponse struct {
	Articles   []BookmarkArticleItem `json:"articles"`
	Pagination PaginationInfo        `json:"pagination"`
}

type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

// ArchiveItem 文章归档项
type ArchiveItem struct {
	YearMonth string `json:"year_month"` // 格式: "2024年11月"
	Count     int    `json:"count"`      // 该月的文章数量
}
