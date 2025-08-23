package message

import "blog/model_def"

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

type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}
