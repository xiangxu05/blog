package model_def

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int32     `gorm:"primaryKey;autoIncrement"`
	Nickname  string    `gorm:"size:32;not null"`              // 添加index加速用户名查询
	Username  string    `gorm:"size:32;not null;unique;index"` // 添加index加速用户名查询
	Password  string    `gorm:"size:64;not null"`              // 增大以支持加密哈希
	Role      string    `gorm:"size:32;default:'user'"`        // 角色字段，默认普通用户
	Email     string    `gorm:"size:64;unique;index"`          // 添加unique和index
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Avatar    string    `gorm:"size:255"` //
}

func (User) TableName() string {
	return "tb_user"
}
