package model_def

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ArticleID int32    `gorm:"index"` // 添加index加速按文章查询
	Article   Article  `gorm:"foreignKey:ArticleID"`
	UserID    int32    `gorm:"index"` // 添加index
	User      int32    `gorm:"foreignKey:UserID"`
	Content   string   `gorm:"type:text;not null"`
	ParentID  *int32   `gorm:"index"` // 添加index支持嵌套评论查询
	Parent    *Comment `gorm:"foreignKey:ParentID"`
}

func (Comment) TableName() string {
	return "tb_comment"
}
