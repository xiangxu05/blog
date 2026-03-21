package message

// CommentCreateRequest 发表评论
type CommentCreateRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *int32 `json:"parent_id"`
}

// CommentAuthor 评论作者摘要
type CommentAuthor struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CommentReplyTo 回复对象（父评论作者）
type CommentReplyTo struct {
	UserID   int32  `json:"user_id"`
	Nickname string `json:"nickname"`
}

// CommentItem 单条评论
type CommentItem struct {
	ID        uint            `json:"id"`
	Content   string          `json:"content"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Edited    bool            `json:"edited"` // updated_at 晚于 created_at，表示曾编辑
	ParentID  *int32          `json:"parent_id"`
	Author    CommentAuthor   `json:"author"`
	ReplyTo   *CommentReplyTo `json:"reply_to,omitempty"`
	LikeCount int             `json:"like_count"`
	Liked     bool            `json:"liked"`
}

// CommentListResponse 评论列表
type CommentListResponse struct {
	Comments   []CommentItem  `json:"comments"`
	Pagination PaginationInfo `json:"pagination"`
}

// CommentCreated 发表评论成功
type CommentCreated struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
