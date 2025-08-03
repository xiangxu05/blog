package model_def

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID         int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	UserID     int32          `gorm:"not null;index" json:"user_id"` // 添加index加速按作者查询
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title      string         `gorm:"size:100;not null;index" json:"title"` // 添加index支持标题搜索
	ArticleID  int32          `gorm:"not null;index" json:"article_id"`
	Content    ArticleContent `gorm:"foreignKey:ArticleID;not null" json:"content,omitempty"`
	Version    int            `gorm:"default:0" json:"version"`
	CategoryID int32          `gorm:"index" json:"category_id"` // 添加index加速分类查询
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Views      int            `gorm:"default:0;index" json:"views"` // 添加index支持热门排序
	Tags       []Tag          `gorm:"many2many:tb_article_tags;" json:"tags,omitempty"`
}

type ArticleContent struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ArticleID   int32          `gorm:"not null;index" json:"article_id"`
	Description string         `gorm:"size:255" json:"description"`
	Version     int            `gorm:"default:0;index" json:"version"`        // 添加index支持版本查询
	Content     string         `gorm:"type:longtext;not null" json:"content"` // 内容使用数组存储，支持版本控制
}

type Category struct {
	ID          int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"size:32;not null;unique;index" json:"name"` // 添加unique和index
	Description string         `gorm:"size:255" json:"description"`
}

type Tag struct {
	ID        int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"size:32;not null;unique;index" json:"name"` // 添加unique和index加速标签查询
}

type ArticleTag struct {
	ArticleID int32 `gorm:"primaryKey;index"` // 添加index
	TagID     int32 `gorm:"primaryKey;index"` // 添加index
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Article) TableName() string {
	return "tb_article"
}

func (ArticleContent) TableName() string {
	return "tb_article_content"
}

func (Category) TableName() string {
	return "tb_category"
}

func (Tag) TableName() string {
	return "tb_tag"
}

func (ArticleTag) TableName() string {
	return "tb_article_tag"
}
