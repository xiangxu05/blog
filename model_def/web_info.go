package model_def

import "time"

type WebInfo struct {
	ID         uint      `gorm:"primaryKey"`
	ArticleNum int       `gorm:"column:article_num"`
	FileNum    int       `gorm:"column:file_num"`
	UserNum    int       `gorm:"column:user_num"`
	Views      int       `gorm:"column:views"`
	Categories int       `gorm:"column:categories"`
	Users      int       `gorm:"column:users"`
	Comments   int       `gorm:"column:comments"`
	StoreUsage int       `gorm:"column:store_usage"`
	Capacity   int       `gorm:"column:capacity"`
	LastUpdate time.Time `gorm:"column:last_update"`
	LastLogin  time.Time `gorm:"column:last_login"`
	LastClear  time.Time `gorm:"column:last_clear"`
	LastBackup time.Time `gorm:"column:last_backup"`
}

func (WebInfo) TableName() string {
	return "web_info"
}
