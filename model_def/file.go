package model_def

import "time"

type File struct {
	ID        int32     `gorm:"primaryKey;autoIncrement"`
	Filename  string    `json:"filename"`
	StorePath string    `json:"store_path"`
	OwnerID   int32     `json:"owner_id"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (File) TableName() string {
	return "tb_file"
}
