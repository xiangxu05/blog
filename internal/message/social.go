package message

// ArticleSocial 文章点赞/收藏状态（用于详情页）
type ArticleSocial struct {
	LikeCount  int  `json:"like_count"`
	Liked      bool `json:"liked"`
	Bookmarked bool `json:"bookmarked"`
}

// CommentUpdateRequest 修改评论
type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required"`
}
