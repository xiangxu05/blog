package model_def

import "gorm.io/gorm"

// ArticleLike 文章点赞（登录用户）
type ArticleLike struct {
	gorm.Model
	UserID    int32 `gorm:"uniqueIndex:uid_aid_like;index;not null"`
	ArticleID int32 `gorm:"uniqueIndex:uid_aid_like;index;not null"`
}

func (ArticleLike) TableName() string { return "tb_article_like" }

// ArticleBookmark 文章收藏
type ArticleBookmark struct {
	gorm.Model
	UserID    int32 `gorm:"uniqueIndex:uid_aid_bm;index;not null"`
	ArticleID int32 `gorm:"uniqueIndex:uid_aid_bm;index;not null"`
}

func (ArticleBookmark) TableName() string { return "tb_article_bookmark" }

// CommentLike 评论点赞
type CommentLike struct {
	gorm.Model
	UserID    int32  `gorm:"uniqueIndex:uid_cid_lk;not null"`
	CommentID uint   `gorm:"uniqueIndex:uid_cid_lk;not null"`
}

func (CommentLike) TableName() string { return "tb_comment_like" }
