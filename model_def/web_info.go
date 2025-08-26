package model_def

import "time"

type WebInfo struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ArticleNum    int       `gorm:"column:article_num" json:"article_num"`
	FileNum       int       `gorm:"column:file_num" json:"file_num"`
	UserNum       int       `gorm:"column:user_num" json:"user_num"`
	CategoriesNum int       `gorm:"column:categories" json:"categories_num"`
	CommentsNum   int       `gorm:"column:comments" json:"comments_num"`
	StoreUsage    string    `gorm:"column:store_usage" json:"store_usage"`
	Capacity      string    `gorm:"column:capacity" json:"capacity"`
	Views         int       `gorm:"column:views" json:"views"`
	LastUpdate    time.Time `gorm:"column:last_update" json:"last_update"`
	LastLogin     time.Time `gorm:"column:last_login" json:"last_login"`
	LastClear     time.Time `gorm:"column:last_clear" json:"last_clear"`
	LastBackup    time.Time `gorm:"column:last_backup" json:"last_backup"`
}

func (WebInfo) TableName() string {
	return "web_info"
}
