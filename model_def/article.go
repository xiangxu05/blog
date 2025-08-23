package model_def

import (
	"time"
)

type Article struct {
	ID        int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 自动维护
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 自动维护
	UserID    int32     `gorm:"not null;index" json:"user_id"`    // 按作者查
	Title     string    `gorm:"size:100;not null;index" json:"title"`
	Version   int       `gorm:"default:0" json:"version"`
	Category  string    `gorm:"size:32;default:'default';index" json:"category"`
	Tags      string    `gorm:"default:''" json:"tags"`       // JSON 存储
	Views     int       `gorm:"default:0;index" json:"views"` // 热门排序
}

func (Article) TableName() string {
	return "tb_article"
}

type ArticleVersion struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ArticleID   int32     `gorm:"not null;index" json:"article_id"`
	Description string    `gorm:"size:255" json:"description"`
	Version     int       `gorm:"default:0;index" json:"version"` // 版本查询
	StoreID     int32     `gorm:"not null;index" json:"store_id"`
}

func (ArticleVersion) TableName() string {
	return "tb_article_content"
}

// todo 加一个这个表
type Category struct {
	ID        int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Name      string    `gorm:"size:32;not null" json:"name"`
	Number    int       `gorm:"default:0" json:"number"`
}

func (Category) TableName() string {
	return "tb_category"
}
