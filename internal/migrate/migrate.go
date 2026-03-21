package migrate

import (
	"blog/model_def"

	"gorm.io/gorm"
)

// AutoMigrate 启动时补齐新表/新列（SQLite 增量迁移，不影响已有数据）
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model_def.ArticleLike{},
		&model_def.ArticleBookmark{},
		&model_def.CommentLike{},
	)
}
