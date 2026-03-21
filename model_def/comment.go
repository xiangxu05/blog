package model_def

import "gorm.io/gorm"

// Comment 文章评论（支持楼中楼：ParentID 指向父评论）
type Comment struct {
	gorm.Model
	ArticleID int32  `gorm:"index;not null"`
	UserID    int32  `gorm:"index;not null"`
	Content   string `gorm:"type:text;not null"`
	ParentID  *int32 `gorm:"index"`
}

func (Comment) TableName() string {
	return "tb_comment"
}
